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
	"github.com/syzoj/syzoj-ng-go/util"
)

type queueItem struct {
	id        primitive.ObjectID
	problemId primitive.ObjectID
	language  string
	code      string
}

type judger struct {
	fetchLock   chan struct{} // size = 1, effectively a lock, held by FetchTask when it is waiting for tasks
	abortNotify chan struct{} // Notifies FetchTask to abort
	listLock    sync.Mutex    // protects judgingTask
	judgingTask []int
}

type Submission struct {
	Lock   sync.RWMutex
	Id primitive.ObjectID
	Broker *util.Broker
	Done   bool
	Score  float64
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
	srv.submissions = make(map[primitive.ObjectID]*Submission)
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
		item := &queueItem{
			id:        submission.Id,
			problemId: submission.Problem,
			language:  submission.Content.Language,
			code:      submission.Content.Code,
		}
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
	item := &queueItem{
		id:        submission.Id,
		problemId: submission.Problem,
		language:  submission.Content.Language,
		code:      submission.Content.Code,
	}
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

func (c *Core) GetSubmission(submissionId primitive.ObjectID) *Submission {
	c.submissionsLock.Lock()
	submission, ok := c.submissions[submissionId]
	if !ok {
		submission = new(Submission)
		submission.Id = submissionId
		submission.Broker = util.NewBroker()
		c.submissions[submissionId] = submission
		submission.Lock.Lock()
	}
	c.submissionsLock.Unlock()
	if !ok {
		var submissionModel model.Submission
		var err error
		if err = c.mongodb.Collection("submission").FindOne(c.context, bson.D{{"_id", submissionId}}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"result", 1}})).Decode(&submissionModel); err != nil {
			log.WithField("submissionId", submissionId).Error("Failed to load submission: ", err)
		}
		if submissionModel.Result.Status == "Done" {
			submission.Done = true
			submission.Score = submissionModel.Result.Score
		} else {
			submission.Done = false
		}
		submission.Lock.Unlock()
	}
	return submission
}
