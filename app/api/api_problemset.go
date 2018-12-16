package api

import (
	"encoding/json"
	"net/http"

	"github.com/syzoj/syzoj-ng-go/app/judge"
	"github.com/syzoj/syzoj-ng-go/app/problemset"

	"github.com/google/uuid"

	"github.com/syzoj/syzoj-ng-go/app/session"
)

type ProblemsetCreateRequest struct{}
type ProblemsetCreateResponse struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
}

func (srv *ApiServer) HandleCreateProblemset(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, r, err)
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
	if id, err = srv.problemsetService.NewProblemset(sess.AuthUserId); err != nil {
		return
	}
	writeResponseWithSession(w, ProblemsetCreateResponse{
		ProblemsetId: id,
	}, sess)
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
			writeError(w, r, err)
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

	if err = srv.problemsetService.AddProblem(req.ProblemsetId, sess.AuthUserId, req.Name, req.ProblemId); err != nil {
		return
	}
	writeResponseWithSession(w, ProblemsetAddProblemResponse{}, sess)
}

type ProblemsetListProblemRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
}
type ProblemsetListProblemResponse struct {
	Problems []ProblemsetListProblemEntry `json:"problems"`
}
type ProblemsetListProblemEntry struct {
	Name  string `json:"name"`
	Title string `json:"title"`
}

func (srv *ApiServer) HandleProblemsetList(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, r, err)
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

	var req ProblemsetListProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	var problems []problemset.ProblemInfo
	if problems, err = srv.problemsetService.ListProblem(req.ProblemsetId, sess.AuthUserId); err != nil {
		return
	}
	var entries []ProblemsetListProblemEntry = make([]ProblemsetListProblemEntry, len(problems))
	for i, problem := range problems {
		entries[i] = ProblemsetListProblemEntry{
			Name:  problem.Name,
			Title: problem.Title,
		}
	}
	writeResponseWithSession(w, ProblemsetListProblemResponse{Problems: entries}, sess)
}

type ProblemsetViewProblemRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
	Name         string    `json:"name"`
}
type ProblemsetViewProblemResponse struct {
	Statement judge.ProblemStatement `json:"statement"`
}

func (srv *ApiServer) HandleProblemsetView(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, r, err)
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

	var req ProblemsetViewProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	var p1 problemset.ProblemInfo
	if p1, err = srv.problemsetService.ViewProblem(req.ProblemsetId, sess.AuthUserId, req.Name); err != nil {
		return
	}
	var p2 *judge.Problem
	if p2, err = srv.judgeService.GetProblem(p1.ProblemId); err != nil {
		return
	}

	writeResponseWithSession(w, ProblemsetViewProblemResponse{
		Statement: p2.Statement,
	}, sess)
}

type ProblemsetSubmitRequest struct {
	ProblemsetId uuid.UUID `json:"problemset_id"`
	ProblemName  string    `json:"problem_name"`
	Type         string    `json:"type"`
	Traditional  judge.TraditionalSubmission
}
type ProblemsetSubmitResponse struct {
	SubmissionId uuid.UUID `json:"submission_id"`
}

func (srv *ApiServer) HandleProblemsetSubmit(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, r, err)
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
		if submissionId, err = srv.problemsetService.SubmitTraditional(req.ProblemsetId, sess.AuthUserId, req.ProblemName, req.Traditional); err != nil {
			return
		}
		writeResponseWithSession(w, ProblemsetSubmitResponse{
			SubmissionId: submissionId,
		}, sess)
		break
	default:
		err = BadRequestError
		return
	}
}
