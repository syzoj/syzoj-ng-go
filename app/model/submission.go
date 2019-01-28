package model

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Submission struct {
	Id               primitive.ObjectID   `bson:"_id"`
	Type             string               `bson:"type"`
	User             primitive.ObjectID   `bson:"user"`
	Owner            []primitive.ObjectID `bson:"owner"`
	Problem          primitive.ObjectID   `bson:"problem"`
	SubmitTime       time.Time            `bson:"submit_time"`
	Content          SubmissionContent    `bson:"content"`
	Result           SubmissionResult     `bson:"result"`
	JudgeQueueStatus JudgeQueueStatus     `bson:"judge_queue_status"`
}

type SubmissionContent struct {
	Language string `bson:"language"`
	Code     string `bson:"code"`
}

type SubmissionResult struct {
	Status      string        `bson:"status"`
	Score       float64       `bson:"score"`
	MemoryUsage int64         `bson:"memory_usage"`
	TimeUsage   time.Duration `bson:"time_usage"`
}

type JudgeQueueStatus struct {
	Version string `bson:"version"`
}
