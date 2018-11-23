package problemset

import "errors"

type ProblemsetPrivilege int
type ProblemsetUserRole interface{}
type ProblemsetInfo interface {
	GetDefaultRole() ProblemsetUserRole
	GetGuestRole() ProblemsetUserRole
	GetRegisteredUserRole() ProblemsetUserRole
	GetCreatorRole() ProblemsetUserRole
	CheckPrivilege(u ProblemsetUserRole, p ProblemsetPrivilege) error
}
type ProblemsetProvider interface {
	GetDefaultProblemsetInfo() ProblemsetInfo
}

const (
	ProblemsetCreateProblemPrivilege = iota
	ProblemsetViewProblemPrivilege
	ProblemsetSubmitProblemPrivilege
)

var ProblemsetPermissionDeniedError = errors.New("Permission denied")
var ProblemsetLoginRequiredError = errors.New("Please login first")
var ProblemsetPermissionInvalidError = errors.New("Invalid permission id")

var problemsetProviders = map[string]ProblemsetProvider{
	"standard": standardProblemsetProvider{},
}

func GetProblemsetType(ptype string) ProblemsetProvider {
	if provider, ok := problemsetProviders[ptype]; ok {
		return provider
	} else {
		return nil
	}
}
