package session

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	NewSession() (uuid.UUID, *Session, error)
	GetSession(uuid.UUID) (*Session, error)
	UpdateSession(uuid.UUID, *Session) error
	GarbageCollect() error
	Close() error
}

type Session struct {
	Version    uuid.UUID
	AuthUserId uuid.UUID
	UserName   string
	Expiry     time.Time
}

var ErrConcurrentUpdate = errors.New("Concurrent update to session")
var ErrSessionNotFound = errors.New("Session not found")
