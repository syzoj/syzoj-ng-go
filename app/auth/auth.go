package auth

import (
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	RegisterUser(userName string, password string) (uuid.UUID, error)
	LoginUser(userName string, password string) (uuid.UUID, error)
	Close() error
}

var ErrInvalidUserName = errors.New("Invalid user name")
var ErrDuplicateUserName = errors.New("Duplicate user name")
var ErrUserNotFound = errors.New("User not found")
var ErrPasswordIncorrect = errors.New("Incorrect password")
var ErrInvalidAuthInfo = errors.New("Invalid auth info")
