package problemset

import (
	"errors"

	"github.com/google/uuid"

    "github.com/syzoj/syzoj-ng-go/app/judge"
)

type Service interface {
	// Creates a new problemset.
	NewProblemset(OwnerId uuid.UUID) (uuid.UUID, error)
	// Adds a traditional problem to the problemset.
	AddProblem(id uuid.UUID, userId uuid.UUID, name string, problemId uuid.UUID) error
	// Views the specified problem.
	ViewProblem(id uuid.UUID, userId uuid.UUID, name string) (ProblemInfo, error)
	// Gets a list of all problems in the problemset.
	ListProblem(id uuid.UUID, userId uuid.UUID) ([]ProblemInfo, error)
	// Submits to a traditional problem.
	SubmitTraditional(id uuid.UUID, userId uuid.UUID, name string, data judge.TraditionalSubmission) (uuid.UUID, error)
	// Views the specified submission.
	ViewSubmission(id uuid.UUID, userId uuid.UUID, submissionId uuid.UUID) (SubmissionInfo, error)
	Close() error
}

type ProblemInfo struct {
	// The name of problem.
	Name string `json:"name"`
	// The title of problem.
	Title string `json:"title"`
	// The problem ID.
	ProblemId uuid.UUID `json:"problem_id"`
}

type SubmissionInfo struct {
	// The type of submission.
	Type string `json:"type"`
    // The submitter.
    UserId uuid.UUID `json:"user"`
    // The name of problem.
    ProblemName string `json:"problem_name"`
    // Whether the problem has been judged.
    Complete bool `json:"complete"`
    // If Type == "traditional", the submitted data.
    Traditional *judge.TraditionalSubmission
    // The result.
    Result judge.TaskCompleteInfo `json:"result"`
}

var ErrInvalidProblemsetType = errors.New("Invalid problemset type")
var ErrProblemsetNotFound = errors.New("Problemset not found")
var ErrOperationNotSupported = errors.New("Operation not supported")
var ErrDuplicateProblemName = errors.New("Duplicate problem name")
var ErrInvalidProblemName = errors.New("Invalid problem name")
var ErrDuplicateUUID = errors.New("UUID dupication")
var ErrPermissionDenied = errors.New("Permission denied")
var ErrAnonymousSubmission = errors.New("Anonymous submission")
var ErrProblemNotFound = errors.New("Problem not found")
var ErrNotImplemented = errors.New("Not implemented")
var ErrSubmissionNotFound = errors.New("Submission not found")
