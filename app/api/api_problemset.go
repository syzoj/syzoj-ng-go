package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/syzoj/syzoj-ng-go/app/judge"
	"github.com/syzoj/syzoj-ng-go/app/problemset"

	"github.com/google/uuid"

	"github.com/syzoj/syzoj-ng-go/app/session"
)

type ProblemsetCreateRequest struct{}
type ProblemsetCreateResponse struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
}

func (srv *ApiServer) HandleCreateProblemset(w http.ResponseWriter, r *http.Request, _ uuid.UUID, sess *session.Session) ApiError {
	var err error
	if err = requireLogin(sess); err != nil {
		return err.(ApiError)
	}
	var id uuid.UUID
	if id, err = srv.problemsetService.NewProblemset(sess.AuthUserId); err != nil {
		return problemsetError(err)
	}
	writeResponse(w, ProblemsetCreateResponse{
		ProblemsetId: id,
	}, sess)
	return nil
}

type ProblemsetAddProblemRequest struct {
	Name      string    `json:"name"`
	ProblemId uuid.UUID `json:"problem_id"`
}
type ProblemsetAddProblemResponse struct{}

func (srv *ApiServer) HandleProblemsetAdd(w http.ResponseWriter, r *http.Request, _ uuid.UUID, sess *session.Session) ApiError {
	var err error
	if err = requireLogin(sess); err != nil {
		return err.(ApiError)
	}
	vars := mux.Vars(r)
	var problemsetId = uuid.MustParse(vars["problemset_id"])
	var req ProblemsetAddProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return problemsetError(err)
	}
	if err = srv.problemsetService.AddProblem(problemsetId, sess.AuthUserId, req.Name, req.ProblemId); err != nil {
		return problemsetError(err)
	}
	writeResponse(w, ProblemsetAddProblemResponse{}, sess)
	return nil
}

type ProblemsetListProblemResponse struct {
	Problems []ProblemsetListProblemEntry `json:"problems"`
}
type ProblemsetListProblemEntry struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func (srv *ApiServer) HandleProblemsetList(w http.ResponseWriter, r *http.Request, _ uuid.UUID, sess *session.Session) ApiError {
	var err error
	if err = requireLogin(sess); err != nil {
		return err.(ApiError)
	}
	vars := mux.Vars(r)
	var problemsetId = uuid.MustParse(vars["problemset_id"])
	var problems []problemset.ProblemInfo
	if problems, err = srv.problemsetService.ListProblem(problemsetId, sess.AuthUserId); err != nil {
		return problemsetError(err)
	}
	var entries []ProblemsetListProblemEntry = make([]ProblemsetListProblemEntry, len(problems))
	for i, problem := range problems {
		entries[i] = ProblemsetListProblemEntry{
			Name:  problem.Name,
			Title: problem.Title,
		}
	}
	writeResponse(w, ProblemsetListProblemResponse{Problems: entries}, sess)
	return nil
}

type ProblemsetViewProblemRequest struct {
	Name string `json:"name"`
}
type ProblemsetViewProblemResponse struct {
	Statement string `json:"statement"`
}

func (srv *ApiServer) HandleProblemsetView(w http.ResponseWriter, r *http.Request, _ uuid.UUID, sess *session.Session) ApiError {
	var err error
	if err = requireLogin(sess); err != nil {
		return err.(ApiError)
	}
	vars := mux.Vars(r)
	var problemsetId = uuid.MustParse(vars["problemset_id"])
	var req ProblemsetViewProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return problemsetError(err)
	}

	var p1 problemset.ProblemInfo
	if p1, err = srv.problemsetService.ViewProblem(problemsetId, sess.AuthUserId, req.Name); err != nil {
		return problemsetError(err)
	}
	var p2 = new(judge.Problem)
	if err = srv.judgeService.GetProblemStatementInfo(p1.ProblemId, p2); err != nil {
		return judgeError(err)
	}
	writeResponse(w, ProblemsetViewProblemResponse{
		Statement: p2.Statement,
	}, sess)
	return nil
}

type ProblemsetSubmitRequest struct {
	ProblemName string `json:"problem_name"`
	Type        string `json:"type"`
	Traditional judge.TraditionalSubmission
}
type ProblemsetSubmitResponse struct {
	SubmissionId uuid.UUID `json:"submission_id"`
}

func (srv *ApiServer) HandleProblemsetSubmit(w http.ResponseWriter, r *http.Request, _ uuid.UUID, sess *session.Session) ApiError {
	var err error
	if err = requireLogin(sess); err != nil {
		return err.(ApiError)
	}
	vars := mux.Vars(r)
	var problemsetId = uuid.MustParse(vars["problemset_id"])
	var req ProblemsetSubmitRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError(err)
	}
	switch req.Type {
	case "traditional":
		var submissionId uuid.UUID
		if submissionId, err = srv.problemsetService.SubmitTraditional(problemsetId, sess.AuthUserId, req.ProblemName, req.Traditional); err != nil {
			return problemsetError(err)
		}
		writeResponse(w, ProblemsetSubmitResponse{
			SubmissionId: submissionId,
		}, sess)
		return nil
	default:
		return ErrNotImplemented
	}
}
