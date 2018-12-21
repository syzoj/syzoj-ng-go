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

func (s *ApiServer) HandleProblemCreate(w http.ResponseWriter, r *http.Request, _ uuid.UUID, sess *session.Session) ApiError {
	var err error
	if err = requireLogin(sess); err != nil {
		return err.(ApiError)
	}
	var req CreateProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError(err)
	}
	info := judge.Problem{
		Owner: sess.AuthUserId,
		Title: req.Title,
	}
	var problemId uuid.UUID
	if problemId, err = s.judgeService.CreateProblem(&info); err != nil {
		return judgeError(err)
	}
	writeResponse(w, CreateProblemResponse{ProblemId: problemId}, sess)
	return nil
}

type ViewProblemResponse struct {
	Title     string `json:"title"`
	Statement string `json:"statement"`
	IsOwner   bool   `json:"is_owner"`
	Token     string `json:"token"`
}

func (s *ApiServer) HandleProblemView(w http.ResponseWriter, r *http.Request, _ uuid.UUID, sess *session.Session) ApiError {
	var err error
	vars := mux.Vars(r)
	var problemId = uuid.MustParse(vars["problem_id"])
	var info = new(judge.Problem)
	if err = s.judgeService.GetProblemFullInfo(problemId, info); err != nil {
		return judgeError(err)
	}
	var resp ViewProblemResponse
	resp.Statement = info.Statement
	resp.Title = info.Title
	if info.Owner == sess.AuthUserId {
		resp.Token = info.Token
		resp.IsOwner = true
	}
	writeResponse(w, &resp, sess)
	return nil
}

type ResetProblemTokenResponse struct {
	Token string `json:"token"`
}

func (s *ApiServer) HandleResetProblemToken(w http.ResponseWriter, r *http.Request, _ uuid.UUID, sess *session.Session) ApiError {
	var err error
	vars := mux.Vars(r)
	var problemId = uuid.MustParse(vars["problem_id"])
	var info = new(judge.Problem)
	if err = s.judgeService.GetProblemOwnerInfo(problemId, info); err != nil {
		return judgeError(err)
	}
	if info.Owner != sess.AuthUserId {
		return ErrPermissionDenied
	}
	if err = s.judgeService.ResetProblemToken(problemId, info); err != nil {
		return judgeError(err)
	}
	var resp ResetProblemTokenResponse
	resp.Token = info.Token
	writeResponse(w, &resp, sess)
	return nil
}

func (s *ApiServer) HandleProblemUpdate(w http.ResponseWriter, r *http.Request, _ uuid.UUID, sess *session.Session) ApiError {
	var err error
	vars := mux.Vars(r)
	var problemId = uuid.MustParse(vars["problem_id"])
	var info = new(judge.Problem)
	if err = s.judgeService.GetProblemOwnerInfo(problemId, info); err != nil {
		return judgeError(err)
	}
	if info.Owner != sess.AuthUserId {
		return ErrPermissionDenied
	}
	if err = s.judgeService.UpdateProblem(problemId, info); err != nil {
		return judgeError(err)
	}
	writeResponse(w, struct{}{}, sess)
	return nil
}

type ProblemChangeTitleRequest struct {
	Title string `json:"title"`
}

func (s *ApiServer) HandleProblemChangeTitle(w http.ResponseWriter, r *http.Request, _ uuid.UUID, sess *session.Session) ApiError {
	var err error
	vars := mux.Vars(r)
	var problemId = uuid.MustParse(vars["problem_id"])
	var req ProblemChangeTitleRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError(err)
	}

	var info = new(judge.Problem)
	if err = s.judgeService.GetProblemOwnerInfo(problemId, info); err != nil {
		return judgeError(err)
	}
	if info.Owner != sess.AuthUserId {
		return ErrPermissionDenied
	}
	info.Title = req.Title
	if err = s.judgeService.ChangeProblemTitle(problemId, info); err != nil {
		return judgeError(err)
	}
	writeResponse(w, struct{}{}, sess)
	return nil
}
