package core

import (
	"context"
	"sync"
	"sync/atomic"

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

func (item *queueItem) getFields() logrus.Fields {
	return logrus.Fields{
		"id":        EncodeObjectID(item.id),
		"problemId": EncodeObjectID(item.problemId),
	}
}

func (srv *Core) initJudge(ctx context.Context) (err error) {
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
	i := atomic.AddInt64(&srv.queueSize, 1)
	srv.queueItems.Store(i, item)
	srv.queue <- i
}

// Notifies that a submission's status has changed.
func (srv *Core) NotifySubmission(id primitive.ObjectID) {
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
