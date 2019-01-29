package model

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Problem struct {
	Id         primitive.ObjectID   `bson:"_id"`
	Title      string               `bson:"title"`
	Statement  string               `bson:"statement,omitempty"`
	Owner      []primitive.ObjectID `bson:"owner"`
	CreateTime time.Time            `bson:"create_time"`
    
    // Couple with the public problemset
    PublicStats ProblemStats `bson:"public_stats"`
}

type ProblemStats struct {
    Accept int64 `bson:"accept"`
    Submission int64 `bson:"submission"`
}
