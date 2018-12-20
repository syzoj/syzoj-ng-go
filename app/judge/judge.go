package judge

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

// The interface for traditional judge service.
type Service interface {
	// Creates a new problem. Only Title and Owner is used; other fields are ignored.
	CreateProblem(info *Problem) (uuid.UUID, error)
	// Refresh the problem and read statement from disk. info is ignored.
	UpdateProblem(id uuid.UUID, info *Problem) error
	// Resets the token for problem.
	ResetProblemToken(id uuid.UUID, info *Problem) (error)
	// Gets the title, statement, token and owner for a problem.
	GetProblemFullInfo(id uuid.UUID, info *Problem) (error)
	// Gets the statement for a problem.
	GetProblemStatementInfo(id uuid.UUID, info *Problem) (error)
	// Gets the owner for a problem.
	GetProblemOwnerInfo(id uuid.UUID, info *Problem) (error)
	// Deletes a problem.
	DeleteProblem(id uuid.UUID) error
	// Adds a submission to queue and receive callback.
	QueueSubmission(sub *Submission, callback Callback) (Task, error)
	// Handles WebSocket connections from judgers.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	// Stops the service.
	Close() error
}

// The type represents all information about a problem.
type Problem struct {
	Title string
	Statement string
	Token     string  
	Owner     uuid.UUID  
}

// The interface represents a task in queue.
type Task interface {
	// Aborts the task.
	Abort() error
}

type Callback interface {
	// Called when a judger picks up the task and starts running.
	OnStart(TaskStartInfo)
	// Called when the task makes new progress.
	OnProgress(TaskProgressInfo)
	// Called when the task completes.
	OnComplete(TaskCompleteInfo)
	// Called if the task is aborted for any reason.
	// Either HandleJudgeComplete() or HandleJudgeError() is called.
	OnError(error)
}

type TaskStartInfo struct{}

type TaskCompleteInfo struct {
	Status string      `json:"status"`
	Score  float64     `json:"score"`
	Detail interface{} `json:"detail"`
}

type TaskProgressInfo struct{}

type Submission struct {
	ProblemId   uuid.UUID
	Traditional TraditionalSubmission
}
type TraditionalSubmission struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

var ErrNotImplemented = errors.New("Not implemented")
var ErrQueueFull = errors.New("Submission queue full")
var ErrConcurrentUpdate = errors.New("Concurrent update")
var ErrProblemNotExist = errors.New("Problem does not exist")
