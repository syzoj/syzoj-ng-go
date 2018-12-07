package traditional

import (
	"github.com/google/uuid"
)

type traditionalJudgeService struct {
	ps    problemsetCallbackService
	queue chan TraditionalSubmission
}

func NewTraditionalJudgeService() (TraditionalJudgeService, error) {
	s := &traditionalJudgeService{
		queue: make(chan TraditionalSubmission),
	}
	return s, nil
}

func (ps *traditionalJudgeService) RegisterProblemsetService(s problemsetCallbackService) {
	if ps.ps != nil {
		panic("traditionalJudgeService: RegisterProblemsetService called twice")
	}
	ps.ps = s
}

func (ps *traditionalJudgeService) QueueSubmission(problemsetId uuid.UUID, submissionId uuid.UUID, submissionData *TraditionalSubmission) error {
	go func() {
		ps.ps.InvokeProblemset(problemsetId, &TraditionalSubmissionResultMessage{
			SubmissionId: submissionId,
			Result: TraditionalSubmissionResult{
				Status: "Not supported",
			},
		}, nil)
	}()
	return nil
}
