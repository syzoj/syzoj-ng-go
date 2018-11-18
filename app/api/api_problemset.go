package api

import (
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

	perm, err := srv.CheckGroupPermission(groupName, sess.AuthUserId, model.GroupCreateProblemsetPrivilege)
	if err == GroupNotFoundError {
		srv.NotFound(w, err)
		return
	} else if err != nil {
		srv.InternalServerError(w, err)
		return
	} else if !perm {
		srv.Forbidden(w, err)
	}

	util.GenerateUUID()
	log.Println("problemset created")
}