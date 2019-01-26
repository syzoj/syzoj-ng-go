package judge

import (
	"context"
	"sync"
	"sync/atomic"
    "errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/sirupsen/logrus"

	judge_api "github.com/syzoj/syzoj-ng-go/app/judge/protos"
	"github.com/syzoj/syzoj-ng-go/app/model"
)

var log = logrus.StandardLogger()

type judgeService struct {
	mongodb    *mongo.Database
	queue      chan int64
	queueSize  int64
	queueItems sync.Map
}

type queueItem struct {
    id primitive.ObjectID
    problemId primitive.ObjectID
    language string
    code string
    version string
}

func (item *queueItem) getFields() logrus.Fields {
	return logrus.Fields{
        "id": EncodeObjectID(item.id),
        "problemId": EncodeObjectID(item.problemId),
    }
}

func CreateJudgeService(mongodb *mongo.Client) (Service, error) {
	srv := &judgeService{
		mongodb: mongodb.Database("syzoj"),
	}
	err := srv.init(context.Background())
	if err != nil {
		return srv, err
	}
	return srv, nil
}

func (srv *judgeService) init(ctx context.Context) (err error) {
	srv.queue = make(chan int64, 1000)
	srv.queueItems = sync.Map{}
	srv.queueSize = 0
    var cursor mongo.Cursor
    if cursor, err = srv.mongodb.Collection("submission").Find(ctx,
        bson.D{{"judge_queue_status", bson.D{{"$exists", true}}}},
        mongo_options.Find().SetProjection(bson.D{{"_id", 1}, {"problem", 1}, {"content.language", 1}, {"content.code", 1}, {"judge_queue_status", 1}})); err != nil {
        return
    }
    defer cursor.Close(ctx)
    for cursor.Next(ctx) {
        submission := new(model.Submission)
        if err = cursor.Decode(&submission); err != nil {
            panic(err)
        }
        go srv.enqueueModel(submission)
    }
    if err = cursor.Err(); err != nil {
        return
    }
	return
}

func (srv *judgeService) Close() error {
	return nil
}

func (srv *judgeService) enqueueModel(model *model.Submission) {
    item := &queueItem{
        id: model.Id,
        problemId: model.Problem,
        language: model.Content.Language,
        code: model.Content.Code,
        version: model.JudgeQueueStatus.Version,
    }
    srv.enqueue(item)
}

func (srv *judgeService) enqueue(item *queueItem) {
	log.WithFields(item.getFields()).Info("Adding submission to queue")
	i := atomic.AddInt64(&srv.queueSize, 1)
	srv.queueItems.Store(i, item)
	srv.queue <- i
}

func (srv *judgeService) NotifySubmission(ctx context.Context, id primitive.ObjectID) (err error) {
    submission := new(model.Submission)
    if err = srv.mongodb.Collection("submission").FindOne(ctx,
        bson.D{{"_id", id}, {"judge_queue_status", bson.D{{"$exists", true}}}},
        mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"problem", 1}, {"content", 1}, {"judge_queue_status", 1}})).Decode(&submission); err != nil {
        return
    }
    go srv.enqueueModel(submission)
    return
}

func (srv *judgeService) FetchTask(ctx context.Context, req *judge_api.JudgeRequest) (res *judge_api.FetchTaskResult, err error) {
    select {
    case id := <- srv.queue:
        itemObj, _ := srv.queueItems.Load(id)
        item := itemObj.(*queueItem)
        res = new(judge_api.FetchTaskResult)
        res.Success = true
        res.Task = &judge_api.Task{
            TaskTag: id,
            ProblemId: EncodeObjectID(item.problemId),
            Language: item.language,
            Code: item.code,
        }
        log.WithFields(item.getFields()).WithField("judgerId", req.GetJudgerId()).Info("Judge item taken by judger")
        return
    default:
        res = new(judge_api.FetchTaskResult)
        res.Success = false
        return
    }
}

func (srv *judgeService) SetTaskProgress(s judge_api.Judge_SetTaskProgressServer) (err error) {
	return
}

func (srv *judgeService) SetTaskResult(ctx context.Context, in *judge_api.TaskResult) (e *empty.Empty, err error) {
    id := in.TaskTag
    itemObj, found := srv.queueItems.Load(id)
    if !found {
        err = errors.New("Invalid taskTag")
        return
    }
    item := itemObj.(*queueItem)
    srv.queueItems.Delete(id)
    var result *mongo.UpdateResult
    if result, err = srv.mongodb.Collection("submission").UpdateOne(ctx,
        bson.D{{"_id", item.id}, {"judge_queue_status.version", item.version}},
        bson.D{
            {"$set", bson.D{
                {"result", bson.D{
                    {"status", in.Result},
                    {"score", in.Score}},
                }},
            },
            {"$unset", bson.D{{"judge_queue_status", 1}}},
        }); err != nil {
        panic(err)
    }
    if result.MatchedCount == 0 {
        // Taken by another judge process
        log.WithFields(item.getFields()).Warning("Conflict judger detected")
    }
    e = new(empty.Empty)
	return
}
