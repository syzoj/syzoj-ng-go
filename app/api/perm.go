package api

import (
	"database/sql"
	"encoding/json"

	model_group "github.com/syzoj/syzoj-ng-go/app/model/group"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

func GetGroupId(cxt *ApiContext) *ApiError {
	row := cxt.tx.QueryRow("SELECT id FROM groups WHERE name=$1", cxt.groupName)
	var groupIdBytes []byte
	if err := row.Scan(&groupIdBytes); err != nil {
		if err == sql.ErrNoRows {
			return GroupNotFoundError
		}
		panic(err)
	}
	if groupId, err := util.UUIDFromBytes(groupIdBytes); err != nil {
		panic(err)
	} else {
		cxt.groupId = groupId
	}
	return nil
}

func GetGroupProblemsetId(cxt *ApiContext) *ApiError {
	row := cxt.tx.QueryRow("SELECT groups.id, problemsets.id FROM groups JOIN problemsets ON groups.id=problemsets.group_id WHERE groups.name=$1 AND problemsets.name=$2", cxt.groupName, cxt.problemsetName)
	var groupIdBytes []byte
	var problemsetIdBytes []byte
	if err := row.Scan(&groupIdBytes, &problemsetIdBytes); err != nil {
		if err == sql.ErrNoRows {
			return ProblemsetNotFoundError
		}
		panic(err)
	}
	if groupId, err := util.UUIDFromBytes(groupIdBytes); err != nil {
		panic(err)
	} else {
		cxt.groupId = groupId
	}
	if problemsetId, err := util.UUIDFromBytes(problemsetIdBytes); err != nil {
		panic(err)
	} else {
		cxt.problemsetId = problemsetId
	}
	return nil
}

func GetGroupProblemsetProblemId(cxt *ApiContext) *ApiError {
	row := cxt.tx.QueryRow("SELECT groups.id, problemsets.id, problem.id FROM groups JOIN problemsets ON groups.id=problemsets.group_id JOIN problems on problemsets.id=problems.problemset_id WHERE groups.name=$1 AND problemsets.name=$2 AND problems.name=$3", cxt.groupName, cxt.problemsetName, cxt.problemName)
	var groupIdBytes []byte
	var problemsetIdBytes []byte
	var problemIdBytes []byte
	if err := row.Scan(&groupIdBytes, &problemsetIdBytes, &problemIdBytes); err != nil {
		if err == sql.ErrNoRows {
			return ProblemNotFoundError
		}
		panic(err)
	}
	if groupId, err := util.UUIDFromBytes(groupIdBytes); err != nil {
		panic(err)
	} else {
		cxt.groupId = groupId
	}
	if problemsetId, err := util.UUIDFromBytes(problemsetIdBytes); err != nil {
		panic(err)
	} else {
		cxt.problemsetId = problemsetId
	}
	if problemId, err := util.UUIDFromBytes(problemIdBytes); err != nil {
		panic(err)
	} else {
		cxt.problemId = problemId
	}
	return nil
}

func GetGroupPolicy(cxt *ApiContext) *ApiError {
	row := cxt.tx.QueryRow("SELECT policy_info FROM groups WHERE id=$1", cxt.groupId.ToBytes())
	var policyInfoBytes []byte
	if err := row.Scan(&policyInfoBytes); err != nil {
		return GroupNotFoundError
	}
	groupProvider := model_group.GetGroupType()
	groupPolicy := groupProvider.GetDefaultGroupPolicy()
	if err := json.Unmarshal(policyInfoBytes, &groupPolicy); err != nil {
		panic(err)
	}
	cxt.groupPolicy = groupPolicy
	return nil
}

func GetGroupUserRole(cxt *ApiContext) *ApiError {
	if !cxt.sess.IsLoggedIn() {
		cxt.groupUserRole = cxt.groupPolicy.GetGuestRole()
	} else {
		row := cxt.tx.QueryRow("SELECT role_info FROM group_users WHERE group_id=$1 AND user_id=$2", cxt.groupId, cxt.sess.AuthUserId)
		var roleInfoBytes []byte
		if err := row.Scan(&roleInfoBytes); err != nil {
			cxt.groupUserRole = cxt.groupPolicy.GetRegisteredUserRole()
			return nil
		}
		groupUserRole := cxt.groupPolicy.GetDefaultRole()
		if err := json.Unmarshal(roleInfoBytes, &groupUserRole); err != nil {
			panic(err)
		}
		cxt.groupUserRole = groupUserRole
	}
	return nil
}

func CheckGroupPrivilege(cxt *ApiContext, priv model_group.GroupPrivilege) *ApiError {
	if cxt.groupPolicy.CheckPrivilege(cxt.groupUserRole, priv) != nil {
		return PermissionDeniedError
	} else {
		return nil
	}
}
