package core

import (
	"errors"
)

var ErrInvalidProblem = errors.New("Invalid problem")
var ErrProblemNotExist = errors.New("Problem doesn't exist")
var ErrConflict = errors.New("Conflict operation")
var ErrDuplicateUserName = errors.New("Duplicate user name")
var ErrInvalidUserName = errors.New("Invalid user name")
var ErrGeneral = errors.New("General failure")
var ErrContestNotRunning = errors.New("Contest not running")
