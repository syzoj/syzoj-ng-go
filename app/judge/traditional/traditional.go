package traditional

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

type problemsetCallbackService interface {
	InvokeProblemset(id uuid.UUID, req interface{}) (interface{}, error)
}

type TraditionalJudgeService interface {
	RegisterProblemsetService(problemsetCallbackService)
	QueueSubmission(problemsetId uuid.UUID, submissionId uuid.UUID, submissionData *TraditionalSubmission) error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Close() error
}

type TraditionalSubmission struct {
	Language  string
	Code      string
	ProblemId uuid.UUID
}

type TraditionalSubmissionResultMessage struct {
	SubmissionId uuid.UUID
	Result       TraditionalSubmissionResult
}

type TraditionalSubmissionResult struct {
	Status string `json:"status"`
}

type TraditionalJudgeMessage struct {
	Tag       int64     `json:"tag"`
	ProblemId uuid.UUID `json:"problem_id"`
	Language  string    `json:"language"`
	Code      string    `json:"code"`
}

type TraditionalJudgeResponse struct {
	Tag    int64                       `json:"tag"`
	Result TraditionalSubmissionResult `json:"result"`
}

var ErrQueueFull = errors.New("Submission queue full")
