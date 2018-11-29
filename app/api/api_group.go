package api

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	model_group "github.com/syzoj/syzoj-ng-go/app/model/group"
)

type GroupCreateRequest struct {
	GroupName string `json:"group_name"`
}

func HandleGroupCreate(cxt *ApiContext) ApiResponse {
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
	groupId := uuid.New()
	groupProvider := model_group.GetGroupType()
	groupPolicy := groupProvider.GetDefaultGroupPolicy()
	_, err := cxt.tx.Exec("INSERT INTO groups (id, name, policy_info) VALUES ($1, $2, $3)", groupId[:], req.GroupName, marshalJson(groupPolicy))
	if err != nil {
		if sqlErr, ok := err.(*pq.Error); ok {
			if sqlErr.Code == "23505" && sqlErr.Constraint == "groups_name_unique" {
				return DuplicateGroupNameError
			}
		}
		panic(err)
	}

	groupCreatorRole := groupPolicy.GetCreatorRole()
	_, err = cxt.tx.Exec("INSERT INTO group_users (group_id, user_id, role_info) VALUES ($1, $2, $3)", groupId[:], cxt.sess.AuthUserId[:], marshalJson(groupCreatorRole))
	if err != nil {
		panic(err)
	}

	DoneTx(cxt)
	return Success(nil)
}
