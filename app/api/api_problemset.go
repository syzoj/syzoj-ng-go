package api

import (
	"encoding/json"
	"github.com/syzoj/syzoj-ng-go/app/model"
	"github.com/gorilla/mux"
	"github.com/syzoj/syzoj-ng-go/app/util"
	"net/http"
	"log"
)

type CreateProblemsetRequest struct {
	ProblemsetName string `json:"problemset_name"`
}
type CreateProblemsetResponse struct {
	Success bool `json:"success"`
	Reason string `json:"reason"`
}
func (srv *ApiServer) HandleProblemsetCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["group-name"]
	sess := srv.GetSession(r)
	if !sess.LoggedIn {
		srv.Forbidden(w, NotLoggedInError)
		return
	}
	reqDecoder := json.NewDecoder(r.Body)
	var req CreateProblemsetRequest
	if err := reqDecoder.Decode(&req); err != nil {
		srv.BadRequest(w, err)
		return
	}

	groupId, perm, err := srv.CheckGroupPermission(groupName, sess.AuthUserId, model.GroupCreateProblemsetPrivilege)
	if err == GroupNotFoundError {
		srv.NotFound(w, err)
		return
	} else if err != nil {
		srv.InternalServerError(w, err)
		return
	} else if !perm {
		srv.Forbidden(w, err)
		return
	}

	problemsetId, err := util.GenerateUUID()
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}
	_, err = srv.db.Query("INSERT INTO problemsets (id, name, group_id, type, info) VALUES ($1, $2, $3, $4, $5)", problemsetId, req.ProblemsetName, groupId.ToBytes(), 1, "{}")
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}
	log.Println("problemset created")
}