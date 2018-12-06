package session

import (
	"errors"
	"github.com/google/uuid"
)

type SessionService interface {
	NewSession() (uuid.UUID, *Session, error)
	GetSession(uuid.UUID) (*Session, error)
	UpdateSession(uuid.UUID, *Session) error
}

type Session struct {
	Version    uuid.UUID
	AuthUserId uuid.UUID
}

var ConcurrentUpdateError = errors.New("Concurrent update to session")