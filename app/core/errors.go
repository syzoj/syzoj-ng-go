package core

import (
	"errors"
)

var ErrInvalidProblem = errors.New("Invalid problem")
var ErrProblemNotExist = errors.New("Problem doesn't exist")
var ErrConflict = errors.New("Conflict operation")
var ErrDuplicateUserName = errors.New("Duplicate user name")
var ErrInvalidUserName = errors.New("Invalid user name")
var ErrDuplicateEmail = errors.New("Duplicate email")
var ErrInvalidEmail = errors.New("Invalid email")
var ErrGeneral = errors.New("General failure")
var ErrContestNotRunning = errors.New("Contest not running")
var ErrAlreadyRegistered = errors.New("Already registered in contest")
var ErrTooManySubmissions = errors.New("Too many submissions")
