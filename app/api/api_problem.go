package api

import (
	_ "github.com/syzoj/syzoj-ng-go/app/model/problem"
	model_problemset "github.com/syzoj/syzoj-ng-go/app/model/problemset"
	"log"
)

type CreateProblemRequest struct {
	GroupName      string `json:"group_name"`
	ProblemsetName string `json:"problemset_name"`
	ProblemName    string `json:"problem_name"`
	ProblemType    string `json:"problem_type"`
}

func HandleProblemCreate(cxt *ApiContext) ApiResponse {
	log.Println("HandleProblemCreate")
	var req CreateProblemRequest
	if err := cxt.ReadBody(&req); err != nil {
		return err
	}
	UseTx(cxt)
	cxt.groupName = req.GroupName
	cxt.problemsetName = req.ProblemsetName
	if err := GetGroupProblemsetId(cxt); err != nil {
		return err
	}
	if err := GetGroupPolicy(cxt); err != nil {
		return err
	}
	if err := GetGroupUserRole(cxt); err != nil {
		return err
	}
	if err := GetProblemsetInfo(cxt); err != nil {
		return err
	}
	if err := GetProblemsetUserRole(cxt); err != nil {
		return err
	}
	if !CheckProblemsetPrivilege(cxt, model_problemset.ProblemsetCreateProblemPrivilege) {
		return PermissionDeniedError
	}
	log.Println("Done")
	return Success("permission pass")
}
