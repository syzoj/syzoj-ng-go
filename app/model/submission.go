package model

import (
    "time"

    "github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Submission struct {
    Id primitive.ObjectID `bson:"_id"`
    Type string `bson:"type"`
    User primitive.ObjectID `bson:"user"`
    Owner []primitive.ObjectID `bson:"owner"`
    Problem primitive.ObjectID `bson:"problem"`
    Language string `bson:"language"`
    Code string `bson:"code"`
    Status string `bson:"status"`
    Score float64 `bson:"score"`
    SubmitTime time.Time `bson:"submit_time"`
}
