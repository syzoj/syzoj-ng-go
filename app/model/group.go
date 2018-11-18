package model

import (
	"errors"
)

type GroupPolicyInfo struct {
	MemberCreateProblemset bool `json:"member_create_problemset"`
}

type UserRoleInfo struct {
	Role int `json:"role"`
}

var GroupOwnerRole = UserRoleInfo{Role: 3}

const (
	GroupCreateProblemsetPrivilege = iota // Create problemsets
	GroupViewProblemsetPrivilege          // View private problemsets
	GroupManageProblemsetPrivilege        // Manage problemsets
	GroupManageGroupPrivilege             // Mange group users
)

var rolePermission = map[int]map[int]bool{
	0: map[int]bool{
		GroupCreateProblemsetPrivilege: false,
	},
}

var InvalidUserRoleInfoError = errors.New("Invalid user role")

func (g *GroupPolicyInfo) CheckPrivilege(u UserRoleInfo, p int) (bool, error) {
	switch u.Role {
	case 0: // Guest
		return false, nil
	case 1: // Member
		switch p {
		case GroupViewProblemsetPrivilege:
			return true, nil
		case GroupCreateProblemsetPrivilege:
			return g.MemberCreateProblemset, nil
		}
		return false, nil
	case 2: // Admin
		switch p {
		case GroupViewProblemsetPrivilege:
			return true, nil
		case GroupCreateProblemsetPrivilege:
			return true, nil
		case GroupManageProblemsetPrivilege:
			return true, nil
		}
		return false, nil
	case 3: // Owner
		return true, nil
	}
	return false, InvalidUserRoleInfoError
}
