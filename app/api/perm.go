package api

import (
	"database/sql"
	"encoding/json"

	"github.com/google/uuid"

	model_group "github.com/syzoj/syzoj-ng-go/app/model/group"
	model_problemset "github.com/syzoj/syzoj-ng-go/app/model/problemset"
)

func GetGroupId(cxt *ApiContext) ApiResponse {
	row := cxt.tx.QueryRow("SELECT id FROM groups WHERE name=$1", cxt.groupName)
	var groupId uuid.UUID
	if err := row.Scan(&groupId); err != nil {
		if err == sql.ErrNoRows {
			return GroupNotFoundError
		}
		panic(err)
	}
	cxt.groupId = groupId
	return nil
}

func GetGroupProblemsetId(cxt *ApiContext) ApiResponse {
	row := cxt.tx.QueryRow("SELECT groups.id, problemsets.id FROM groups JOIN problemsets ON groups.id=problemsets.group_id WHERE groups.name=$1 AND problemsets.name=$2", cxt.groupName, cxt.problemsetName)
	var groupId uuid.UUID
	var problemsetId uuid.UUID
	if err := row.Scan(&groupId, &problemsetId); err != nil {
		if err == sql.ErrNoRows {
			return ProblemsetNotFoundError
		}
		panic(err)
	}
	cxt.groupId = groupId
	cxt.problemsetId = problemsetId
	return nil
}

func GetGroupProblemsetProblemId(cxt *ApiContext) ApiResponse {
	row := cxt.tx.QueryRow("SELECT groups.id, problemsets.id, problem.id FROM groups JOIN problemsets ON groups.id=problemsets.group_id JOIN problems on problemsets.id=problems.problemset_id WHERE groups.name=$1 AND problemsets.name=$2 AND problems.name=$3", cxt.groupName, cxt.problemsetName, cxt.problemName)
	var groupId uuid.UUID
	var problemsetId uuid.UUID
	var problemId uuid.UUID
	if err := row.Scan(&groupId, &problemsetId, &problemId); err != nil {
		if err == sql.ErrNoRows {
			return ProblemNotFoundError
		}
		panic(err)
	}
	cxt.groupId = groupId
	cxt.problemsetId = problemsetId
	cxt.problemId = problemId
	return nil
}

func GetGroupPolicy(cxt *ApiContext) ApiResponse {
	row := cxt.tx.QueryRow("SELECT policy_info FROM groups WHERE id=$1", cxt.groupId[:])
	var policyInfoBytes []byte
	if err := row.Scan(&policyInfoBytes); err != nil {
		if err == sql.ErrNoRows {
			return GroupNotFoundError
		}
		panic(err)
	}
	groupProvider := model_group.GetGroupType()
	groupPolicy := groupProvider.GetDefaultGroupPolicy()
	if err := json.Unmarshal(policyInfoBytes, &groupPolicy); err != nil {
		panic(err)
	}
	cxt.groupPolicy = groupPolicy
	return nil
}

func GetGroupUserRole(cxt *ApiContext) ApiResponse {
	if !cxt.sess.IsLoggedIn() {
		cxt.groupUserRole = cxt.groupPolicy.GetGuestRole()
	} else {
		row := cxt.tx.QueryRow("SELECT role_info FROM group_users WHERE group_id=$1 AND user_id=$2", cxt.groupId[:], cxt.sess.AuthUserId[:])
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

func CheckGroupPrivilege(cxt *ApiContext, priv model_group.GroupPrivilege) bool {
	if cxt.groupPolicy.CheckPrivilege(cxt.groupUserRole, priv) != nil {
		return false
	} else {
		return true
	}
}

func GetProblemsetInfo(cxt *ApiContext) ApiResponse {
	row := cxt.tx.QueryRow("SELECT type, info FROM problemsets WHERE id=$1", cxt.problemsetId[:])
	var problemsetType string
	var infoBytes []byte
	if err := row.Scan(&problemsetType, &infoBytes); err != nil {
		if err == sql.ErrNoRows {
			return ProblemsetNotFoundError
		}
		panic(err)
	}
	problemsetProvider := model_problemset.GetProblemsetType(problemsetType)
	problemsetInfo := problemsetProvider.GetDefaultProblemsetInfo()
	if err := json.Unmarshal(infoBytes, &problemsetInfo); err != nil {
		panic(err)
	}
	cxt.problemsetInfo = problemsetInfo
	return nil
}

func GetProblemsetUserRole(cxt *ApiContext) ApiResponse {
	if !cxt.sess.IsLoggedIn() {
		cxt.problemsetUserRole = cxt.problemsetInfo.GetGuestRole()
	} else {
		row := cxt.tx.QueryRow("SELECT info FROM problemset_users WHERE problemset_id=$1 AND user_id=$2", cxt.problemsetId[:], cxt.sess.AuthUserId[:])
		var infoBytes []byte
		if err := row.Scan(&infoBytes); err != nil {
			cxt.problemsetUserRole = cxt.problemsetInfo.GetRegisteredUserRole()
			return nil
		}
		problemsetUserRole := cxt.problemsetInfo.GetDefaultRole()
		if err := json.Unmarshal(infoBytes, &problemsetUserRole); err != nil {
			panic(err)
		}
		cxt.problemsetUserRole = problemsetUserRole
	}
	return nil
}

func CheckProblemsetPrivilege(cxt *ApiContext, priv model_problemset.ProblemsetPrivilege) bool {
	if cxt.groupPolicy.CheckProblemsetPrivilege(cxt.groupUserRole, priv) == nil {
		return true
	}
	if cxt.problemsetInfo.CheckPrivilege(cxt.problemsetUserRole, priv) == nil {
		return true
	}
	return false
}
