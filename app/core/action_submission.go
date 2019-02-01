package core

import (
	"context"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
)

type Submit1 struct {
	ProblemId primitive.ObjectID
	Submitter primitive.ObjectID
	Enqueue   bool
	Code      Code
}
type Code struct {
	Language string
	Code     string
}
type Submit1Resp struct {
	SubmissionId primitive.ObjectID
}

// Submits to the public problemset.
// Possible errors:
// * ErrProblemNotExist
// * MongoDB error or context error
func (c *Core) Action_Submit(ctx context.Context, req *Submit1) (*Submit1Resp, error) {
	var err error
	if _, err = c.mongodb.Collection("problem").FindOne(ctx, bson.D{{"_id", req.ProblemId}}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}})).DecodeBytes(); err != nil {
		return nil, ErrProblemNotExist
	}
	submissionId := primitive.NewObjectID()
	document := bson.D{
		{"_id", submissionId},
		{"type", "standard"},
		{"user", req.Submitter},
		{"owner", []primitive.ObjectID{req.Submitter}},
		{"problem", req.ProblemId},
		{"content", bson.D{
			{"language", req.Code.Language},
			{"code", req.Code.Code},
		}},
		{"submit_time", time.Now()},
	}
	if req.Enqueue {
		var versionBytes [16]byte
		rand.Read(versionBytes[:])
		document = append(document, bson.E{"judge_queue_status", bson.D{{"version", hex.EncodeToString(versionBytes[:])}}})
	}
	if _, err = c.mongodb.Collection("submission").InsertOne(ctx, document); err != nil {
		return nil, err
	}
	go func() {
		result, err := c.mongodb.Collection("problem").UpdateOne(context.TODO(), bson.D{{"_id", req.ProblemId}}, bson.D{{"$inc", bson.D{{"public_stats.submission", 1}}}})
		if err == nil && result.MatchedCount == 0 {
			err = ErrProblemNotExist
		}
		if err != nil {
			log.WithField("problemId", req.ProblemId).WithField("submissionId", submissionId).Error("Failed to increment submission count by one")
		}
	}()
	if req.Enqueue {
		go c.NotifySubmission(submissionId)
	}
	return &Submit1Resp{SubmissionId: submissionId}, nil
}
