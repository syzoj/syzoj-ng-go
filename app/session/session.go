package session

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type SessionService interface {
	NewSession() (uuid.UUID, *Session, error)
	GetSession(uuid.UUID) (*Session, error)
	UpdateSession(uuid.UUID, *Session) error
	Close() error
}

type Session struct {
	Version    uuid.UUID
	AuthUserId uuid.UUID
	Expiry     time.Time
}

var ConcurrentUpdateError = errors.New("Concurrent update to session")
var ErrSessionNotFound = errors.New("Session not found")
