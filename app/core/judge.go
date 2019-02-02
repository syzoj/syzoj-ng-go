package core

import (
	"context"
	"sync"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type queueItem struct {
	id        primitive.ObjectID
	problemId primitive.ObjectID
	language  string
	code      string
	version   string
}

type judger struct {
	fetchLock   chan struct{} // size = 1, effectively a lock, held by FetchTask when it is waiting for tasks
	abortNotify chan struct{} // Notifies FetchTask to abort
	listLock    sync.Mutex    // protects judgingTask
	judgingTask []int
}

type submissionHandler struct {
	lock sync.Mutex
	subscribers map[SubmissionSubscriber]struct{}
	done bool
	score float64
}

type SubmissionSubscriber interface {
	HandleNewScore(done bool, score float64)
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
	srv.submissionHandlers = make(map[primitive.ObjectID]*submissionHandler)
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
		srv.enqueueModel(submission)
	}
	if err = cursor.Err(); err != nil {
		return
	}
	return
}

func (srv *Core) enqueueModel(model *model.Submission) {
	item := &queueItem{
		id:        model.Id,
		problemId: model.Problem,
		language:  model.Content.Language,
		code:      model.Content.Code,
		version:   model.JudgeQueueStatus.Version,
	}
	srv.enqueue(item)
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
func (srv *Core) EnqueueSubmission(id primitive.ObjectID) {
	var err error
	submission := new(model.Submission)
	if err = srv.mongodb.Collection("submission").FindOne(srv.context,
		bson.D{{"_id", id}, {"judge_queue_status", bson.D{{"$exists", true}}}},
		mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"problem", 1}, {"content", 1}, {"judge_queue_status", 1}})).Decode(&submission); err != nil {
		log.WithField("submissionId", id).Error("NotifySubmission: Failed to find submission: " + err.Error())
		return
	}
	srv.enqueueModel(submission)
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

func (c *Core) loadSubmission(submissionId primitive.ObjectID) *submissionHandler {
	c.submissionHandlersLock.Lock()
	handler, ok := c.submissionHandlers[submissionId]
	if !ok {
		handler = new(submissionHandler)
		handler.subscribers = make(map[SubmissionSubscriber]struct{})
		c.submissionHandlers[submissionId] = handler
		handler.lock.Lock()
	}
	c.submissionHandlersLock.Unlock()
	if !ok {
		var submissionModel model.Submission
		var err error
		if err = c.mongodb.Collection("submission").FindOne(c.context, bson.D{{"_id", submissionId}}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"result", 1}})).Decode(&submissionModel); err != nil {
			log.WithField("submissionId", submissionId).Error("Failed to load submission: ", err)
		}
		if submissionModel.Result.Status == "Done" {
			handler.done = true
			handler.score = submissionModel.Result.Score
		} else {
			handler.done = false
		}
	} else {
		handler.lock.Lock()
	}
	return handler
}

// subscriber must be comparable
func (c *Core) SubscribeSubmission(submissionId primitive.ObjectID, subscriber SubmissionSubscriber) {
	handler := c.loadSubmission(submissionId)
	defer handler.lock.Unlock()
	handler.subscribers[subscriber] = struct{}{}
	go subscriber.HandleNewScore(handler.done, handler.score)
}

func (c *Core) UnsubscribeSubmission(submissionId primitive.ObjectID, subscriber SubmissionSubscriber) {
	handler := c.loadSubmission(submissionId)
	defer handler.lock.Unlock()
	delete(handler.subscribers, subscriber)
	if len(handler.subscribers) == 0 {
		log.WithField("submissionId", submissionId).Debug("TODO: Submission has no subscriber, remove from memory")
	}
}
