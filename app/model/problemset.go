package model

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Problemset struct {
	Id             primitive.ObjectID `bson:"_id"`
	Problems       []ProblemsetEntry  `bsan:"problems,omitempty"`
	ProblemsetName string             `bson:"problemset_name,omitempty"`
	Contest        Contest            `bson:"contest,omitempty"`
}

type ProblemsetEntry struct {
	ProblemId primitive.ObjectID `bson:"problem_id,omitempty"`
}

type Contest struct {
	StartTime time.Time            `bson:"start_time,omitempty"`
	EndTime   time.Time            `bson:"end_time,omitempty"`
	Owner     []primitive.ObjectID `bson:"owner,omitempty"`
	Running   bool                 `bson:"running,omitempty"`
	State     string               `bson:"state,omitempty"`
	Schedule  []ContestSchedule    `bson:"schedule,omitempty"`
}

type ContestSchedule struct {
	Key       string    `bson:"key"`
	Done      bool      `bson:"done"`
	Operation bson.Raw  `bson:"operation"`
	StartTime time.Time `bson:"start_time"`
}
