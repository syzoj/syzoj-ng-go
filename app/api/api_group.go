package api

import (
	"github.com/lib/pq"
	model_group "github.com/syzoj/syzoj-ng-go/app/model/group"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type GroupCreateRequest struct {
	GroupName string `json:"group_name"`
}

func HandleGroupCreate(cxt *ApiContext) *ApiError {
	var req GroupCreateRequest
	if err := cxt.ReadBody(&req); err != nil {
		return err
	}
	if !cxt.sess.IsLoggedIn() {
		return NotLoggedInError
	}
	if err := UseTx(cxt); err != nil {
		return err
	}
	groupId, err := util.GenerateUUID()
	if err != nil {
		panic(err)
	}
	groupProvider := model_group.GetGroupType()
	groupPolicy := groupProvider.GetDefaultGroupPolicy()
	_, err = cxt.tx.Exec("INSERT INTO groups (id, name, policy_info) VALUES ($1, $2, $3)", groupId.ToBytes(), req.GroupName, marshalJson(groupPolicy))
	if err != nil {
		if sqlErr, ok := err.(*pq.Error); ok {
			if sqlErr.Code == "23505" && sqlErr.Constraint == "groups_group_name_unique" {
				return DuplicateGroupNameError
			}
		}
		panic(err)
	}

	groupCreatorRole := groupPolicy.GetCreatorRole()
	_, err = cxt.tx.Exec("INSERT INTO group_users (group_id, user_id, role_info) VALUES ($1, $2, $3)", groupId.ToBytes(), cxt.sess.AuthUserId.ToBytes(), marshalJson(groupCreatorRole))
	if err != nil {
		panic(err)
	}

	DoneTx(cxt)
	return nil
}
