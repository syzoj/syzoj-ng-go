package group

import (
	"errors"

	model_problemset "github.com/syzoj/syzoj-ng-go/app/model/problemset"
)

type GroupPrivilege int
type GroupUserRole interface{}
type GroupPolicy interface {
	GetDefaultRole() GroupUserRole
	GetGuestRole() GroupUserRole
	GetRegisteredUserRole() GroupUserRole
	GetCreatorRole() GroupUserRole
	CheckPrivilege(u GroupUserRole, p GroupPrivilege) error
	CheckProblemsetPrivilege(u GroupUserRole, p model_problemset.ProblemsetPrivilege) error
}
type GroupProvider interface {
	GetDefaultGroupPolicy() GroupPolicy
}

const (
	GroupCreateProblemsetPrivilege = iota // Create problemsets
	GroupViewProblemsetPrivilege          // View private problemsets
	GroupManageProblemsetPrivilege        // Manage problemsets
	GroupManageGroupPrivilege             // Mange group users
)

var GroupPermissionDeniedError = errors.New("Permission denied")

var InvalidUserRoleInfoError = errors.New("Invalid user role")

func GetGroupType() GroupProvider {
	return standardGroupProvider{}
}
