package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lib/pq"
	"github.com/syzoj/syzoj-ng-go/app/model/group"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type GroupCreateRequest struct {
	GroupName string `json:"group_name"`
}
type GroupCreateResponse struct{}

func (srv *ApiServer) HandleGroupCreate(w http.ResponseWriter, r *http.Request) {
	sess := srv.GetSession(r)
	reqDecoder := json.NewDecoder(r.Body)
	var req GroupCreateRequest
	if err := reqDecoder.Decode(&req); err != nil {
		srv.BadRequest(w, err)
		return
	}

	if !sess.IsLoggedIn() {
		srv.Forbidden(w, NotLoggedInError)
		return
	}
	trans, err := srv.db.Begin()
	if err != nil {
		panic(err)
	}
	success := false
	defer func() {
		if success {
			err := trans.Commit()
			if err != nil {
				log.Println("Failed to commit:", err)
			}
		} else {
			err := trans.Rollback()
			if err != nil {
				log.Println("Failed to rollback:", err)
			}
		}
	}()

	groupId, err := util.GenerateUUID()
	if err != nil {
		panic(err)
	}
	groupProvider := group.GetGroupType()
	groupPolicy := groupProvider.GetDefaultGroupPolicy()
	_, err = trans.Exec("INSERT INTO groups (id, group_name, policy_info) VALUES ($1, $2, $3)", groupId.ToBytes(), req.GroupName, marshalJson(groupPolicy))
	if err != nil {
		if sqlErr, ok := err.(*pq.Error); ok {
			if sqlErr.Code == "23505" && sqlErr.Constraint == "groups_group_name_unique" {
				srv.SuccessWithError(w, DuplicateGroupNameError)
				return
			}
		}
		panic(err)
	}

	groupCreatorRole := groupPolicy.GetCreatorRole()
	_, err = trans.Exec("INSERT INTO group_users (group_id, user_id, role_info) VALUES ($1, $2, $3)", groupId.ToBytes(), sess.AuthUserId.ToBytes(), marshalJson(groupCreatorRole))
	if err != nil {
		panic(err)
	}

	success = true
	srv.Success(w, GroupCreateResponse{})
}
