package api

import "errors"

type ApiError struct {
	Code    int
	Message string
}

func (e *ApiError) Error() string {
	return e.Message
}
func (e *ApiError) Execute(cxt *ApiContext) {
	cxt.code = e.Code
	cxt.resp = ErrorResponse{Error: e.Message}
}

var ApiEndpointNotFoundError = &ApiError{404, "API endpoint not found"}
var GroupNotFoundError = &ApiError{404, "Group not found"}
var ProblemsetNotFoundError = &ApiError{404, "Problemset not found"}
var ProblemNotFoundError = &ApiError{404, "Problem not found"}
var PermissionDeniedError = &ApiError{403, "Permission denied"}
var NotLoggedInError = &ApiError{401, "Not logged in"}
var DuplicateGroupNameError = &ApiError{200, "Duplicate group name"}
var DuplicateUserNameError = &ApiError{200, "Duplicate user name"}
var DuplicateProblemsetNameError = &ApiError{200, "Duplicate problemset name"}
var BadRequestError = &ApiError{400, "Bad request"}
var InvalidProblemsetTypeError = &ApiError{400, "Invalid or unsupported problemset type"}
var InvalidProblemTypeError = &ApiError{400, "Invalid or unsupported problem type"}
var InternalServerError = &ApiError{500, "Internal server error"}

var AlreadyLoggedInError = &ApiError{200, "Already logged in"}
var UnknownUsernameError = &ApiError{200, "Unknown username"}
var CannotLoginError = &ApiError{200, "Cannot login yet"}
var TwoFactorNotSupportedError = &ApiError{200, "Two factor auth not supported"}
var PasswordIncorrectError = &ApiError{200, "Password incorrect"}

// Internal error
var InvalidAuthUserIdError = errors.New("Invalid AuthUserId")
