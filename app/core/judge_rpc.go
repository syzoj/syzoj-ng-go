package core

import (
	"context"
	"errors"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"

	judge_api "github.com/syzoj/syzoj-ng-go/app/core/protos"
)

type judgeRpc struct {
	*Core
}

func (srv *Core) JudgeRpc() judge_api.JudgeServer {
	return judgeRpc{srv}
}

func (srv judgeRpc) FetchTask(ctx context.Context, in *judge_api.JudgeRequest) (res *judge_api.FetchTaskResult, err error) {
	var judgerId primitive.ObjectID
	var ok bool
	if judgerId, ok = DecodeObjectIDOK(in.JudgerId); !ok {
		err = errors.New("Invalid judger id")
		return
	}
	judger := srv.getJudger(judgerId)
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
		res = new(judge_api.FetchTaskResult)
		res.Success = true
		res.Task = &judge_api.Task{
			TaskTag:   int64(id),
			ProblemId: EncodeObjectID(item.problemId),
			Language:  item.language,
			Code:      item.code,
		}
		log.WithFields(item.getFields()).WithField("judgerId", judgerId).Debug("Judge item taken by judger")
		return
		/*
			default:
				res = new(judge_api.FetchTaskResult)
				res.Success = false
				return
		*/
	}
}

func (srv judgeRpc) SetTaskProgress(s judge_api.Judge_SetTaskProgressServer) (err error) {
	return
}

func (srv judgeRpc) SetTaskResult(ctx context.Context, in *judge_api.SetTaskResultMessage) (e *empty.Empty, err error) {
	var judgerId primitive.ObjectID
	var ok bool
	if judgerId, ok = DecodeObjectIDOK(in.JudgerId); !ok {
		err = errors.New("Invalid judger id")
		return
	}
	judger := srv.getJudger(judgerId)
	id := int(in.TaskTag)
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
	srv.queueLock.Lock()
	item := srv.queueItems[id]
	delete(srv.queueItems, id)
	srv.queueLock.Unlock()
	go func() {
		handler := srv.loadSubmission(item.id)
		handler.done = true
		handler.score = float64(in.Result.Score)
		for subscriber := range handler.subscribers {
			go subscriber.HandleNewScore(true, float64(in.Result.Score))
		}
		handler.lock.Unlock()
		var result *mongo.UpdateResult
		if result, err = srv.mongodb.Collection("submission").UpdateOne(ctx,
			bson.D{{"_id", item.id}, {"judge_queue_status.version", item.version}},
			bson.D{
				{"$set", bson.D{
					{"result", bson.D{
						{"status", in.Result.Result},
						{"score", in.Result.Score}},
					}},
				},
				{"$unset", bson.D{{"judge_queue_status", 1}}},
			}); err != nil {
			panic(err)
		}
		if result.MatchedCount == 0 {
			// Taken by another judge process
			log.WithFields(item.getFields()).Warning("Failed to update judge status due to conflict")
		}
	}()
	e = new(empty.Empty)
	return
}
