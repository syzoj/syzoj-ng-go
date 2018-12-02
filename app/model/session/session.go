package session

import (
	"github.com/google/uuid"
)

type Session struct {
	AuthUserId uuid.UUID `json:"auth_user_id"`
}
