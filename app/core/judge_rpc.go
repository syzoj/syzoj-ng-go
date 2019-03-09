package core

import (
	"context"
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/syzoj/syzoj-ng-go/app/model"
	"github.com/syzoj/syzoj-ng-go/judger/rpc"
)

type judgeRpc struct {
	*Core
}

func (srv *Core) JudgeRpc() rpc.JudgeServer {
	return judgeRpc{srv}
}

func (srv judgeRpc) FetchTask(ctx context.Context, in *rpc.JudgeRequest) (res *rpc.FetchTaskResult, err error) {
	var judgerId primitive.ObjectID
	if judgerId, err = model.GetObjectID(in.JudgerId); err != nil {
		return
	}
	judger := srv.getJudger(judgerId)
	if !judger.checkToken(in.GetJudgerToken()) {
		err = ErrPermissionDenied
		return
	}
loop:
	for {
		select {
		case judger.fetchLock <- struct{}{}:
			defer func() {
				<-judger.fetchLock
			}()
			break loop
		case judger.abortNotify <- struct{}{}:
		}
	}
	log.WithField("judgerId", judgerId).Debug("Judger fetching tasks")
	for _, task := range judger.judgingTask {
		// TODO: potential deadlock here
		srv.queue <- task
	}
	judger.judgingTask = nil
	select {
	case <-ctx.Done():
		log.WithField("judgerId", judgerId).Debug(ctx.Err())
		return nil, ctx.Err()
	case <-judger.abortNotify:
		log.WithField("judgerId", judgerId).Debug("Aborted")
		return nil, errors.New("Aborted")
	case id := <-srv.queue:
		judger.listLock.Lock()
		judger.judgingTask = append(judger.judgingTask, id)
		judger.listLock.Unlock()
		srv.queueLock.Lock()
		item, ok := srv.queueItems[id]
		srv.queueLock.Unlock()
		if !ok {
			panic("Queue item doesn't exist")
		}
		res = new(rpc.FetchTaskResult)
		res.Success = proto.Bool(true)
		res.Task = &rpc.Task{
			TaskTag:           proto.Int64(int64(id)),
			ProblemId:         model.ObjectIDProto(item.problemId),
			SubmissionContent: item.submissionContent,
			ProblemData:       item.problemData,
			ProblemType:       item.problemType,
		}
		log.WithFields(item.getFields()).WithField("judgerId", judgerId).Debug("Judge item taken by judger")
		return
		/*
			default:
				res = new(rpc.FetchTaskResult)
				res.Success = false
				return
		*/
	}
}

func (c judgeRpc) SetTaskProgress(s rpc.Judge_SetTaskProgressServer) (err error) {
	return
}

func (c judgeRpc) SetTaskResult(ctx context.Context, in *rpc.SetTaskResultMessage) (e *empty.Empty, err error) {
	var judgerId primitive.ObjectID
	if judgerId, err = model.GetObjectID(in.JudgerId); err != nil {
		return
	}
	judger := c.getJudger(judgerId)
	id := int(in.GetTaskTag())
	judger.listLock.Lock()
	var found bool
	for k, v := range judger.judgingTask {
		if v == id {
			judger.judgingTask = append(judger.judgingTask[:k], judger.judgingTask[k+1:]...)
			found = true
		}
	}
	judger.listLock.Unlock()
	if !found {
		err = errors.New("Invalid taskTag")
		return
	}
	c.queueLock.Lock()
	item := c.queueItems[id]
	delete(c.queueItems, id)
	c.queueLock.Unlock()
	go func() {
		var result *mongo.UpdateResult
		if result, err = c.mongodb.Collection("submission").UpdateOne(c.context,
			bson.D{{"_id", item.id}},
			bson.D{
				{"$set", bson.D{
					{"result", in.Result},
				}},
				{"$unset", bson.D{{"judge_queue_status", 1}}},
			}); err != nil {
			log.WithField("submissionId", item.id).Error("Failed to update judge queue status: ", err)
			return
		}
		if result.MatchedCount == 0 {
			log.WithFields(item.getFields()).Warning("Failed to update judge status")
		}
		c.invokeSubmissionHook(item.id, in.Result)
	}()
	e = new(empty.Empty)
	return
}
