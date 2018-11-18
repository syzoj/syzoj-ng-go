package api

import (
	"github.com/syzoj/syzoj-ng-go/app/model"
	"encoding/json"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

func (srv *ApiServer) CheckGroupPermission(groupName string, userId util.UUID, privilege int) (groupId util.UUID, success bool, err error) {
	rows, err := srv.db.Query("SELECT id, policy_info FROM groups WHERE groups.group_name=$1", groupName)
	if err != nil {
		return
	}
	if !rows.Next() {
		err = GroupNotFoundError
		return
	}
	var groupPolicyInfoBytes []byte
	var groupIdBytes []byte
	rows.Scan(&groupIdBytes, &groupPolicyInfoBytes)
	groupId, err = util.UUIDFromBytes(groupIdBytes)
	if err != nil {
		return
	}
	var groupPolicyInfo model.GroupPolicyInfo
	err = json.Unmarshal(groupPolicyInfoBytes, &groupPolicyInfo)
	if err != nil {
		return
	}

	var roleInfo model.UserRoleInfo
	rows, err = srv.db.Query("SELECT role_info FROM group_users WHERE groups.id=$1 AND users.id=$2", groupId.ToBytes(), userId.ToBytes())
	if err != nil {
		return
	}
	if rows.Next() {
		var roleInfoBytes []byte
		rows.Scan(roleInfoBytes)
		json.Unmarshal(roleInfoBytes, &roleInfo)
	}
	success, err = groupPolicyInfo.CheckPrivilege(roleInfo, privilege)
	return
}