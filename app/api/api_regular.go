package api

import (
	"encoding/json"
	"net/http"

	"github.com/syzoj/syzoj-ng-go/app/problemset_regular"

	"github.com/google/uuid"

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

	var id uuid.UUID
	if id, err = srv.psregularService.NewProblemset(sess.AuthUserId); err != nil {
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

	var problemId uuid.UUID
	if problemId, err = uuid.NewRandom(); err != nil {
		return
	}
	if err = srv.psregularService.AddTraditionalProblem(req.ProblemsetId, sess.AuthUserId, req.Name, problemId); err != nil {
		return
	}
	writeResponse(w, AddTraditionalProblemResponse{
		ProblemId: problemId,
	})
}

type SubmitTraditionalProblemRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
	ProblemName  string    `json:"problem_name"`
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
	var submissionId uuid.UUID
	if submissionId, err = srv.psregularService.SubmitTraditional(req.ProblemsetId, sess.AuthUserId, req.ProblemName, problemset_regular.TraditionalSubmissionRequest{
		Language: req.Language,
		Code:     req.Code,
	}); err != nil {
		return
	}
	writeResponse(w, SubmitTraditionalProblemResponse{
		SubmissionId: submissionId,
	})
}
