package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/syzoj/syzoj-ng-go/app/problemset"
	"github.com/syzoj/syzoj-ng-go/app/session"
)

type RegularProblemsetCreateRequest struct{}
type RegularProblemsetCreateResponse struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
}

func (srv *ApiServer) HandleRegularProblemsetCreate(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, err)
		}
	}()

	var sess *session.Session
	if _, sess, err = srv.ensureSession(w, r); err != nil {
		return
	}
	if sess.AuthUserId == defaultUserId {
		err = NotLoggedInError
		return
	}

	preq := &problemset.RegularCreateRequest{
		OwnerId: sess.AuthUserId,
	}
	var id uuid.UUID
	if id, err = srv.problemsetService.NewProblemset("regular", preq); err != nil {
		return
	}
	writeResponse(w, RegularProblemsetCreateResponse{
		ProblemsetId: id,
	})
}

type AddTraditionalProblemRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
	Name         string    `json:"name"`
}
type AddTraditionalProblemResponse struct {
	ProblemId uuid.UUID `json:"problem_id"`
}

func (srv *ApiServer) HandleAddTraditionalProblem(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, err)
		}
	}()

	var sess *session.Session
	if _, sess, err = srv.ensureSession(w, r); err != nil {
		return
	}
	if sess.AuthUserId == defaultUserId {
		err = NotLoggedInError
		return
	}

	var req AddTraditionalProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	preq := problemset.RegularAddTraditionalProblemRequest{
		UserId: sess.AuthUserId,
		Name:   req.Name,
	}
	var presp problemset.RegularAddTraditionalProblemResponse
	if err = srv.problemsetService.InvokeProblemset(req.ProblemsetId, &preq, &presp); err != nil {
		return
	}
	writeResponse(w, AddTraditionalProblemResponse{
		ProblemId: presp.ProblemId,
	})
}

type SubmitTraditionalProblemRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
	ProblemId    uuid.UUID `json:"problem_id"`
	Language     string    `json:"language"`
	Code         string    `json:"code"`
}
type SubmitTraditionalProblemResponse struct {
	SubmissionId uuid.UUID `json:"submission_id"`
}

func (srv *ApiServer) HandleSubmitTraditionalProblem(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, err)
		}
	}()

	var sess *session.Session
	if _, sess, err = srv.ensureSession(w, r); err != nil {
		return
	}
	if sess.AuthUserId == defaultUserId {
		err = NotLoggedInError
		return
	}

	var req SubmitTraditionalProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	preq := problemset.RegularSubmitTraditionalProblemRequest{
		ProblemId: req.ProblemId,
		UserId:    sess.AuthUserId,
		Language:  req.Language,
		Code:      req.Code,
	}
	var presp problemset.RegularSubmitTraditionalProblemResponse
	if err = srv.problemsetService.InvokeProblemset(req.ProblemsetId, &preq, &presp); err != nil {
		return
	}
	writeResponse(w, SubmitTraditionalProblemResponse{
		SubmissionId: presp.SubmissionId,
	})
}
