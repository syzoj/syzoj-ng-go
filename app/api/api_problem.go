package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/syzoj/syzoj-ng-go/app/judge"
	"github.com/syzoj/syzoj-ng-go/app/session"
)

type CreateProblemRequest struct {
	Title string `json:"title"`
}
type CreateProblemResponse struct {
	ProblemId uuid.UUID `json:"problem_id"`
}

func (s *ApiServer) HandleProblemCreate(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, r, err)
		}
	}()

	var sess *session.Session
	if _, sess, err = s.ensureSession(w, r); err != nil {
		return
	}
	if sess.AuthUserId == defaultUserId {
		err = NotLoggedInError
		return
	}

	var req CreateProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	info := judge.Problem{
		Owner: sess.AuthUserId,
		Title: req.Title,
	}
	var problemId uuid.UUID
	if problemId, err = s.judgeService.CreateProblem(&info); err != nil {
		return
	}
	writeResponseWithSession(w, CreateProblemResponse{ProblemId: problemId}, sess)
}

type ViewProblemResponse struct {
	Title     string `json:"title"`
	Statement string `json:"statement"`
	IsOwner   bool   `json:"is_owner"`
	Token     string `json:"token"`
}

func (s *ApiServer) HandleProblemView(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, r, err)
		}
	}()

	var sess *session.Session
	if _, sess, err = s.ensureSession(w, r); err != nil {
		return
	}

	vars := mux.Vars(r)
	var problemId uuid.UUID
	if problemId, err = uuid.Parse(vars["problem_id"]); err != nil {
		return
	}
	var info = new(judge.Problem)
	if err = s.judgeService.GetProblemFullInfo(problemId, info); err != nil {
		return
	}

	var resp ViewProblemResponse
	resp.Statement = info.Statement
	resp.Title = info.Title
	if info.Owner == sess.AuthUserId {
		resp.Token = info.Token
		resp.IsOwner = true
	}
	writeResponseWithSession(w, &resp, sess)
}

type ResetProblemTokenResponse struct {
	Token string `json:"token"`
}

func (s *ApiServer) HandleResetProblemToken(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, r, err)
		}
	}()

	var sess *session.Session
	if _, sess, err = s.ensureSession(w, r); err != nil {
		return
	}

	vars := mux.Vars(r)
	var problemId uuid.UUID
	if problemId, err = uuid.Parse(vars["problem_id"]); err != nil {
		return
	}
	var info = new(judge.Problem)
	if err = s.judgeService.GetProblemOwnerInfo(problemId, info); err != nil {
		return
	}
	if info.Owner != sess.AuthUserId {
		err = PermissionDeniedError
		return
	}
	if err = s.judgeService.ResetProblemToken(problemId, info); err != nil {
		return
	}

	var resp ResetProblemTokenResponse
	resp.Token = info.Token
	writeResponseWithSession(w, &resp, sess)
}

func (s *ApiServer) HandleProblemUpdate(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, r, err)
		}
	}()

	var sess *session.Session
	if _, sess, err = s.ensureSession(w, r); err != nil {
		return
	}

	vars := mux.Vars(r)
	var problemId uuid.UUID
	if problemId, err = uuid.Parse(vars["problem_id"]); err != nil {
		return
	}
	var info = new(judge.Problem)
	if err = s.judgeService.GetProblemOwnerInfo(problemId, info); err != nil {
		return
	}
	if info.Owner != sess.AuthUserId {
		err = PermissionDeniedError
		return
	}
	if err = s.judgeService.UpdateProblem(problemId, info); err != nil {
		return
	}

	writeResponseWithSession(w, struct{}{}, sess)
}
