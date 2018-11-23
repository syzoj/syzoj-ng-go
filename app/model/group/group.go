package group

import (
	"errors"
)

type GroupPrivilege int
type GroupUserRole interface{}
type GroupPolicy interface {
	GetDefaultRole() GroupUserRole
	GetGuestRole() GroupUserRole
	GetRegisteredUserRole() GroupUserRole
	GetCreatorRole() GroupUserRole
	CheckPrivilege(u GroupUserRole, p GroupPrivilege) error
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
var GroupPermissionInvalidError = errors.New("Invalid permission id")

func GetGroupType() GroupProvider {
	return standardGroupProvider{}
}
