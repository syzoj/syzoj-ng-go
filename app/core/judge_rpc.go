package core

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"

	judge_api "github.com/syzoj/syzoj-ng-go/app/core/protos"
)

type judgeRpc struct {
	*Core
}

func (srv *Core) JudgeRpc() judge_api.JudgeServer {
	return judgeRpc{srv}
}

func (srv judgeRpc) FetchTask(ctx context.Context, req *judge_api.JudgeRequest) (res *judge_api.FetchTaskResult, err error) {
	select {
	case id := <-srv.queue:
		itemObj, _ := srv.queueItems.Load(id)
		item := itemObj.(*queueItem)
		res = new(judge_api.FetchTaskResult)
		res.Success = true
		res.Task = &judge_api.Task{
			TaskTag:   id,
			ProblemId: EncodeObjectID(item.problemId),
			Language:  item.language,
			Code:      item.code,
		}
		log.WithFields(item.getFields()).WithField("judgerId", req.GetJudgerId()).Info("Judge item taken by judger")
		return
	default:
		res = new(judge_api.FetchTaskResult)
		res.Success = false
		return
	}
}

func (srv judgeRpc) SetTaskProgress(s judge_api.Judge_SetTaskProgressServer) (err error) {
	return
}

func (srv judgeRpc) SetTaskResult(ctx context.Context, in *judge_api.TaskResult) (e *empty.Empty, err error) {
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
