package git

import (
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type GitService interface {
	// Serves HTTP request to /git.
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	// Creates a new git repository.
	CreateRepository(HookType string) (uuid.UUID, error)
	// Resets the token and returns the new token.
	ResetToken(id uuid.UUID) (string, error)
	// Attaches a hook handler.
	AttachHookHandler(HookType string, Handler GitHookHandler)
	Close() error
}

type GitObject [20]byte

type GitHookHandler interface{}

type GitPreReceiveHookHandler interface {
	PreReceiveHook(id uuid.UUID, entries []GitEntry, writer io.Writer) bool
}
type GitUpdateHookHandler interface {
	UpdateHook(id uuid.UUID, entry GitEntry, writer io.Writer) bool
}
type GitPostReceiveHookHandler interface {
	PostReceiveHook(id uuid.UUID, entries []GitEntry, writer io.Writer)
}

type GitEntry struct {
	OldRev  GitObject
	NewRev  GitObject
	RefName string
}

var ErrRepoNotFound = errors.New("Repository not found")
