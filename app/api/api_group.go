package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/lib/pq"
	"github.com/syzoj/syzoj-ng-go/app/model"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type GroupCreateRequest struct {
	GroupName string `json:"group_name"`
}
type GroupCreateResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}

func (srv *ApiServer) HandleGroupCreate(w http.ResponseWriter, r *http.Request) {
	jsonDecoder := json.NewDecoder(r.Body)
	var req GroupCreateRequest
	if err := jsonDecoder.Decode(&req); err != nil {
		srv.BadRequest(w, err)
		return
	}

	session := srv.GetSession(r)
	if !session.LoggedIn {
		srv.Success(w, GroupCreateResponse{Success: false, Reason: "Not logged in"})
		return
	}

	trans, err := srv.db.Begin()
	if err != nil {
		srv.InternalServerError(w, err)
		return
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

	group_id, err := util.GenerateUUID()
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}
	_, err = trans.Exec("INSERT INTO groups (id, group_name, policy_info) VALUES ($1, $2, '{}'::jsonb)", group_id.ToBytes(), req.GroupName)
	if err != nil {
		if sqlErr, ok := err.(*pq.Error); ok {
			if sqlErr.Code == "23505" && sqlErr.Constraint == "groups_group_name" {
				srv.Success(w, GroupCreateResponse{Success: false, Reason: "Duplicate group name"})
				return
			}
		}
		srv.InternalServerError(w, err)
		return
	}

	roleInfoBytes, err := json.Marshal(model.GroupOwnerRole)
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}
	_, err = trans.Exec("INSERT INTO group_users (group_id, user_id, role_info) VALUES ($1, $2, $3)", group_id.ToBytes(), session.AuthUserId.ToBytes(), roleInfoBytes)
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}
	success = true
	srv.Success(w, GroupCreateResponse{Success: true})
}
