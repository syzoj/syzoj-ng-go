package model

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Contest struct {
	Id                   primitive.ObjectID `bson:"_id"`
	Problems             []*ProblemEntry    `bson:"problems,omitempty"`
	Name                 string             `bson:"name,omitempty"`
	Description          string             `bson:"description,omitempty"`
	Owner                primitive.ObjectID `bson:"owner,omitempty"`
	StartTime            time.Time          `bson:"start_time,omitempty"`
	EndTime              time.Time          `bson:"end_time,omitempty"`
	Running              bool               `bson:"running,omitempty"`
	Schedule             []ContestSchedule  `bson:"schedule,omitempty"`
	RanklistType         string             `bson:"ranklist_type",omitempty"`
	RanklistComp         string             `bson:"ranklist_comp",omitempty"`
	JudgeInContest       bool               `bson:"judge_in_contest"`
	SubmissionPerProblem int32              `bson:"submission_per_problem"`
}

type ProblemEntry struct {
	Name      string             `bson:"name,omitempty"`
	ProblemId primitive.ObjectID `bson:"problem_id,omitempty"`
}

type ContestSchedule struct {
	Type      string    `bson:"type"`
	Done      bool      `bson:"done"`
	Operation bson.Raw  `bson:"operation"`
	StartTime time.Time `bson:"start_time"`
}
