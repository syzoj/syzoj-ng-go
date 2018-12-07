package traditional

import (
	"github.com/google/uuid"
)

type problemsetCallbackService interface {
	InvokeProblemset(id uuid.UUID, req interface{}, resp interface{}) error
}

type TraditionalJudgeService interface {
	RegisterProblemsetService(problemsetCallbackService)
	QueueSubmission(problemsetId uuid.UUID, submissionId uuid.UUID, submissionData *TraditionalSubmission) error
}

type TraditionalSubmission struct {
	Language string
	Code     string
}

type TraditionalSubmissionResult struct {
	Status string
}

type TraditionalSubmissionResultMessage struct {
	SubmissionId uuid.UUID
	Result       TraditionalSubmissionResult
}
