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

func (srv *ApiServer) HandleCreateProblemset(w http.ResponseWriter, r *http.Request) {
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

type ProblemsetAddProblemRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
	Name         string    `json:"name"`
	ProblemId    uuid.UUID `json:"problem_id"`
}
type ProblemsetAddProblemResponse struct{}

func (srv *ApiServer) HandleProblemsetAdd(w http.ResponseWriter, r *http.Request) {
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

	var req ProblemsetAddProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	if err = srv.psregularService.AddProblem(req.ProblemsetId, sess.AuthUserId, req.Name, req.ProblemId); err != nil {
		return
	}
	writeResponse(w, ProblemsetAddProblemResponse{})
}

type ProblemsetSubmitRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
	ProblemName  string    `json:"problem_name"`
	Type         string    `json:"type"`
	Traditional  *problemset_regular.TraditionalSubmissionRequest
}
type ProblemsetSubmitResponse struct {
	SubmissionId uuid.UUID `json:"submission_id"`
}

func (srv *ApiServer) HandleProblemsetSubmit(w http.ResponseWriter, r *http.Request) {
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

	var req ProblemsetSubmitRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	switch req.Type {
	case "traditional":
		var submissionId uuid.UUID
		if submissionId, err = srv.psregularService.SubmitTraditional(req.ProblemsetId, sess.AuthUserId, req.ProblemName, *req.Traditional); err != nil {
			return
		}
		writeResponse(w, ProblemsetSubmitResponse{
			SubmissionId: submissionId,
		})
		break
	default:
		err = BadRequestError
		return
	}
}
