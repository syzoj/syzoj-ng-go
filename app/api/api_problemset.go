package api

import (
	"github.com/google/uuid"
	"github.com/lib/pq"

	model_group "github.com/syzoj/syzoj-ng-go/app/model/group"
	model_problemset "github.com/syzoj/syzoj-ng-go/app/model/problemset"
)

type CreateProblemsetRequest struct {
	GroupName      string `json:"group_name"`
	ProblemsetName string `json:"problemset_name"`
	ProblemsetType string `json:"problemset_type"`
}

func HandleProblemsetCreate(cxt *ApiContext) ApiResponse {
	var req CreateProblemsetRequest
	if err := cxt.ReadBody(&req); err != nil {
		return err
	}
	if !cxt.sess.IsLoggedIn() {
		return NotLoggedInError
	}
	UseTx(cxt)
	cxt.groupName = req.GroupName
	if err := GetGroupId(cxt); err != nil {
		return err
	}
	if err := GetGroupPolicy(cxt); err != nil {
		return err
	}
	if err := GetGroupUserRole(cxt); err != nil {
		return err
	}
	if !CheckGroupPrivilege(cxt, model_group.GroupCreateProblemsetPrivilege) {
		return PermissionDeniedError
	}

	problemsetId := uuid.New()
	problemsetProvider := model_problemset.GetProblemsetType(req.ProblemsetType)
	if problemsetProvider == nil {
		return InvalidProblemsetTypeError
	}
	problemsetInfo := problemsetProvider.GetDefaultProblemsetInfo()
	_, err := cxt.tx.Exec(
		"INSERT INTO problemsets (id, name, group_id, type, info) VALUES ($1, $2, $3, $4, $5)",
		problemsetId[:],
		req.ProblemsetName,
		cxt.groupId[:],
		"standard",
		marshalJson(problemsetInfo),
	)
	if err != nil {
		if sqlErr, ok := err.(*pq.Error); ok {
			if sqlErr.Constraint == "problemsets_name_unique" {
				return DuplicateProblemsetNameError
			}
		}
		panic(err)
	}

	problemsetUserRole := problemsetInfo.GetCreatorRole()
	_, err = cxt.tx.Exec(
		"INSERT INTO problemset_users (problemset_id, user_id, info) VALUES ($1, $2, $3)",
		problemsetId[:],
		cxt.sess.AuthUserId[:],
		marshalJson(problemsetUserRole),
	)
	if err != nil {
		panic(err)
	}
	DoneTx(cxt)
	return Success(nil)
}
