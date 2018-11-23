package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lib/pq"

	model_group "github.com/syzoj/syzoj-ng-go/app/model/group"
	model_problemset "github.com/syzoj/syzoj-ng-go/app/model/problemset"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type CreateProblemsetRequest struct {
	ProblemsetName string `json:"problemset_name"`
	ProblemsetType string `json:"problemset_type"`
}
type CreateProblemsetResponse struct{}

func (srv *ApiServer) HandleProblemsetCreate(w http.ResponseWriter, r *http.Request) {
	sess := srv.GetSession(r)
	vars := mux.Vars(r)
	groupName := vars["group-name"]
	reqDecoder := json.NewDecoder(r.Body)
	var req CreateProblemsetRequest
	if err := reqDecoder.Decode(&req); err != nil {
		srv.BadRequest(w, err)
		return
	}

	groupId, groupPolicy := srv.GetGroupPolicyByName(groupName)
	if groupPolicy == nil {
		srv.NotFound(w, GroupNotFoundError)
		return
	}
	userRole := srv.GetGroupUserRole(groupId, groupPolicy, sess.AuthUserId)
	if err := groupPolicy.CheckPrivilege(userRole, model_group.GroupCreateProblemsetPrivilege); err != nil {
		srv.Forbidden(w, CreateProblemsetDeniedError)
		return
	}

	problemsetId, err := util.GenerateUUID()
	if err != nil {
		panic(err)
	}
	problemsetProvider := model_problemset.GetProblemsetType(req.ProblemsetType)
	if problemsetProvider == nil {
		srv.BadRequest(w, InvalidProblemsetTypeError)
		return
	}
	problemsetInfo := problemsetProvider.GetDefaultProblemsetInfo()
	_, err = srv.db.Exec(
		"INSERT INTO problemsets (id, name, group_id, type, info) VALUES ($1, $2, $3, $4, $5)",
		problemsetId.ToBytes(),
		req.ProblemsetName,
		groupId.ToBytes(),
		"standard",
		marshalJson(problemsetInfo),
	)
	if err != nil {
		if sqlErr, ok := err.(*pq.Error); ok {
			if sqlErr.Constraint == "problemsets_name_unique" {
				srv.SuccessWithError(w, DuplicateProblemsetNameError)
				return
			}
		}
		panic(err)
	}
	srv.Success(w, CreateProblemsetResponse{})
}
