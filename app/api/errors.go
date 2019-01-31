package api

type ApiError interface {
	Code() int
	Error() string
}

type apiError struct {
	code    int
	message string
}

func (e *apiError) Code() int {
	return e.code
}

func (e *apiError) Error() string {
	return e.message
}

type internalServerErrorType struct {
	Err error
}

func (e internalServerErrorType) Code() int {
	return 500
}

func (e internalServerErrorType) Error() string {
	return "Internal server error"
}

func internalServerError(err error) ApiError {
	return internalServerErrorType{err}
}

type badRequestErrorType struct {
	Err error
}

func (e badRequestErrorType) Code() int {
	return 400
}

func (e badRequestErrorType) Error() string {
	return e.Err.Error()
}

func badRequestError(err error) ApiError {
	return badRequestErrorType{err}
}

var ErrConflict = &apiError{503, "Please retry"}
var ErrNotImplemented = &apiError{501, "Not implemented"}
var ErrNotFound = &apiError{404, "Not found"}

var ErrProblemNotFound = &apiError{404, "Problem not found"}
var ErrSubmissionNotFound = &apiError{404, "Submission not found"}
var ErrQueueFull = &apiError{503, "Submission queue full"}

var ErrAlreadyLoggedIn = &apiError{200, "Already logged in"}
var ErrNotLoggedIn = &apiError{401, "Authentication required"}
var ErrPermissionDenied = &apiError{403, "Permission denied"}

var ErrUserNotFound = &apiError{200, "User not found"}
var ErrInvalidUserName = &apiError{400, "Invalid user name"}
var ErrDuplicateUserName = &apiError{200, "Duplicate user name"}
var ErrPasswordIncorrect = &apiError{200, "Password incorrect"}
var ErrCannotLogin = &apiError{403, "Cannot login"}

var ErrDuplicatePublicName = &apiError{200, "Duplicate problem name"}
var ErrInvalidPublicName = &apiError{400, "Invalid public name"}
var ErrProblemAlreadyPublic = &apiError{200, "Problem already public"}

var ErrCSRF = &apiError{403, "CSRF token didn't match"}
