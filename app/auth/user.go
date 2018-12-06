package auth

import (
	"errors"

	"github.com/google/uuid"
)

type AuthService interface {
	RegisterUser(userName string, password string) (uuid.UUID, error)
	LoginUser(userName string, password string) (uuid.UUID, error)
}

var ErrDuplicateUserName = errors.New("Duplicate user name")
var ErrPasswordIncorrect = errors.New("Incorrect password")
