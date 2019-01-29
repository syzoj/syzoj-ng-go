package core

import (
    "math/rand"
    "context"
    "time"
    "encoding/hex"

    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/bson/primitive"
    mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
)

type Submit1 struct {
    ProblemId primitive.ObjectID
    Submitter primitive.ObjectID
    Enqueue bool
    Language string
    Code string
}
type Submit1Resp struct {
    SubmissionId primitive.ObjectID
}
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
            {"language", req.Language},
            {"code", req.Code},
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
    if req.Enqueue {
        go c.NotifySubmission(submissionId)
    }
    return &Submit1Resp{SubmissionId: submissionId}, nil
}
