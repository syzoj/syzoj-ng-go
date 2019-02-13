package model

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Contest struct {
	Id          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	Owner       primitive.ObjectID `bson:"owner,omitempty"`
	State       ContestState       `bson:"state,omitempty"`
}

// Volatile contest states.
type ContestState struct {
	Running              bool               `bson:"running,omitempty"`
	StartTime            time.Time          `bson:"start_time,omitempty"`
	Problems             []*ProblemEntry    `bson:"problems,omitempty"`
	Schedule             []*ContestSchedule `bson:"schedule,omitempty"`
	RanklistType         string             `bson:"ranklist_type,omitempty"`
	RanklistComp         string             `bson:"ranklist_comp,omitempty"`
	RanklistVisibility   string             `bson:"ranklist_visibility,omitempty"`
	JudgeInContest       bool               `bson:"judge_in_contest,omitempty"`
	SubmissionPerProblem int32              `bson:"submission_per_problem,omitempty"`
}

type ProblemEntry struct {
	Name      string             `bson:"name,omitempty"`
	ProblemId primitive.ObjectID `bson:"problem_id,omitempty"`
}

type ContestSchedule struct {
	Type      string    `bson:"type"`
	Done      bool      `bson:"done"`
	StartTime time.Time `bson:"start_time"`
}
