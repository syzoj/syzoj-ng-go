package core

import (
	"context"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	// "github.com/mongodb/mongo-go-driver/mongo"
	// mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
)

type ProblemDbNew1 struct {
	Title string
	Owner primitive.ObjectID
}
type ProblemDbNew1Resp struct {
	ProblemId primitive.ObjectID
}

// Possible errors:
// * ErrInvalidProblem
// * MongoDB or context error
func (c *Core) Action_ProblemDb_New(ctx context.Context, req *ProblemDbNew1) (*ProblemDbNew1Resp, error) {
	// TODO: Verify that owner is a real user
	var err error
	if req.Title == "" {
		return nil, ErrInvalidProblem
	}
	problemId := primitive.NewObjectID()
	if _, err = c.mongodb.Collection("problem").InsertOne(ctx, bson.D{
		{"_id", problemId},
		{"title", req.Title},
		{"owner", []primitive.ObjectID{req.Owner}},
		{"create_time", time.Now()},
		{"public_stats", bson.D{
			{"accept", 0},
			{"submission", 0},
		}},
	}); err != nil {
		return nil, err
	}
	return &ProblemDbNew1Resp{ProblemId: problemId}, nil
}
