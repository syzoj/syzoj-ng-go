package auth

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

type AuthService interface {
	RegisterUser(userName string, password string) (uuid.UUID, error)
	LoginUser(userName string, password string) (uuid.UUID, error)
}

var ErrInvalidUserName = errors.New("Invalid user name")
var ErrDuplicateUserName = errors.New("Duplicate user name")
var ErrUserNotFound = errors.New("User not found")
var ErrPasswordIncorrect = errors.New("Incorrect password")

var userNameRegex = regexp.MustCompile("^[0-9A-Za-z]{3,32}$")

func checkUserName(userName string) bool {
	return userNameRegex.MatchString(userName)
}
