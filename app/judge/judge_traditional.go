package judge

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
)

// The interface for traditional judge service.
type Service interface {
	// Creates a new problem.
	CreateProblem(info *Problem) (uuid.UUID, error)
	// Updates a problem.
	UpdateProblem(id uuid.UUID, info *Problem) error
	// Gets a problem.
	GetProblem(id uuid.UUID) (*Problem, error)
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
	Statement ProblemStatement `json:"statement"`
	Owner     uuid.UUID        `json:"owner"`
	Version   int64            `json:"version"`
}
type ProblemStatement struct {
	Title     string `json:"title"`
	Statement string `json:"statement"`
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
	Status string `json:"status"`
}

type TaskProgressInfo struct{}

type Submission struct {
	Language  string
	Code      string
	ProblemId uuid.UUID
}

var ErrNotImplemented = errors.New("Not implemented")
var ErrQueueFull = errors.New("Submission queue full")
var ErrConcurrentUpdate = errors.New("Concurrent update")
var ErrProblemNotExist = errors.New("Problem does not exist")
