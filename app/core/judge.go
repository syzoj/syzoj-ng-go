package core

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type queueItem struct {
	id        primitive.ObjectID
	problemId primitive.ObjectID
	content   *model.SubmissionContent
}

type judger struct {
	fetchLock   chan struct{} // size = 1, effectively a lock, held by FetchTask when it is waiting for tasks
	abortNotify chan struct{} // Notifies FetchTask to abort
	listLock    sync.Mutex    // protects judgingTask
	judgingTask []int
}

type SubmissionHook interface {
	OnSubmissionResult(submissionId primitive.ObjectID, result *model.SubmissionResult)
}

func (c *Core) AddSubmissionHook(hook SubmissionHook) {
	c.submissionHooksMutex.Lock()
	defer c.submissionHooksMutex.Unlock()
	c.submissionHooks[hook] = struct{}{}
}

func (c *Core) RemoveSubmissionHook(hook SubmissionHook) {
	c.submissionHooksMutex.Lock()
	defer c.submissionHooksMutex.Unlock()
	delete(c.submissionHooks, hook)
}

func (c *Core) invokeSubmissionHook(submissionId primitive.ObjectID, result *model.SubmissionResult) {
	c.submissionHooksMutex.Lock()
	defer c.submissionHooksMutex.Unlock()
	for hook := range c.submissionHooks {
		hook.OnSubmissionResult(submissionId, result)
	}
}

func (item *queueItem) getFields() logrus.Fields {
	return logrus.Fields{
		"id":        EncodeObjectID(item.id),
		"problemId": EncodeObjectID(item.problemId),
	}
}

func (srv *Core) initJudge(ctx context.Context) (err error) {
	srv.queue = make(chan int, 1000)
	srv.queueItems = make(map[int]*queueItem)
	srv.queueSize = 0
	srv.judgers = make(map[primitive.ObjectID]*judger)
	var cursor *mongo.Cursor
	if cursor, err = srv.mongodb.Collection("submission").Find(ctx,
		bson.D{{"judge_queue_status.in_queue", true}},
		mongo_options.Find().SetProjection(bson.D{{"_id", 1}, {"problem", 1}, {"content.language", 1}, {"content.code", 1}, {"judge_queue_status", 1}})); err != nil {
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		submission := new(model.Submission)
		if err = cursor.Decode(&submission); err != nil {
			panic(err)
		}
		item := new(queueItem)
		item.content = submission.Content
		item.id = model.MustGetObjectID(submission.GetId())
		item.problemId = model.MustGetObjectID(submission.GetProblem())
		srv.enqueue(item)
	}
	if err = cursor.Err(); err != nil {
		return
	}
	return
}

func (srv *Core) enqueue(item *queueItem) {
	log.WithFields(item.getFields()).Info("Adding submission to queue")
	srv.queueLock.Lock()
	defer srv.queueLock.Unlock()
	i := srv.queueSize
	srv.queueSize++
	srv.queueItems[i] = item
	srv.queue <- i
}

// Puts a submission into queue.
func (c *Core) EnqueueSubmission(id primitive.ObjectID) {
	var err error
	submission := new(model.Submission)
	if err = c.mongodb.Collection("submission").FindOne(c.context,
		bson.D{{"_id", id}},
		mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"problem", 1}, {"content", 1}, {"judge_queue_status", 1}})).Decode(&submission); err != nil {
		log.WithField("submissionId", id).Error("EnqueueSubmission: Failed to find submission: ", err)
		return
	}
	if _, err = c.mongodb.Collection("submission").UpdateOne(c.context,
		bson.D{{"_id", id}},
		bson.D{{"$set", bson.D{{"judge_queue_status", bson.D{{"in_queue", true}}}}}}); err != nil {
		log.WithField("submissionId", id).Error("Failed to set judge queue status: ", err)
		return
	}
	item := new(queueItem)
	item.content = submission.Content
	item.id = model.MustGetObjectID(submission.GetId())
	item.problemId = model.MustGetObjectID(submission.GetProblem())
	c.enqueue(item)
}

func (c *Core) getJudger(id primitive.ObjectID) *judger {
	c.judgerLock.Lock()
	defer c.judgerLock.Unlock()
	j, ok := c.judgers[id]
	if !ok {
		j = new(judger)
		j.abortNotify = make(chan struct{})
		j.fetchLock = make(chan struct{}, 1)
		c.judgers[id] = j
	}
	return j
}
