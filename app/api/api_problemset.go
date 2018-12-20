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

func (srv *ApiServer) HandleCreateProblemset(w http.ResponseWriter, r *http.Request) {
	var err error
	var sess *session.Session
	defer func() {
		if err != nil {
			writeError(w, r, err, sess)
		}
	}()

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
	writeResponse(w, ProblemsetCreateResponse{
		ProblemsetId: id,
	}, sess)
}

type ProblemsetAddProblemRequest struct {
	Name      string    `json:"name"`
	ProblemId uuid.UUID `json:"problem_id"`
}
type ProblemsetAddProblemResponse struct{}

func (srv *ApiServer) HandleProblemsetAdd(w http.ResponseWriter, r *http.Request) {
	var err error
	var sess *session.Session
	defer func() {
		if err != nil {
			writeError(w, r, err, sess)
		}
	}()

	if _, sess, err = srv.ensureSession(w, r); err != nil {
		return
	}
	if sess.AuthUserId == defaultUserId {
		err = NotLoggedInError
		return
	}

	vars := mux.Vars(r)
	var problemsetId uuid.UUID
	if problemsetId, err = uuid.Parse(vars["problemset_id"]); err != nil {
		return
	}
	var req ProblemsetAddProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	if err = srv.problemsetService.AddProblem(problemsetId, sess.AuthUserId, req.Name, req.ProblemId); err != nil {
		return
	}
	writeResponse(w, ProblemsetAddProblemResponse{}, sess)
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
	var sess *session.Session
	defer func() {
		if err != nil {
			writeError(w, r, err, sess)
		}
	}()

	if _, sess, err = srv.ensureSession(w, r); err != nil {
		return
	}
	if sess.AuthUserId == defaultUserId {
		err = NotLoggedInError
		return
	}

	vars := mux.Vars(r)
	var problemsetId uuid.UUID
	if problemsetId, err = uuid.Parse(vars["problemset_id"]); err != nil {
		return
	}

	var problems []problemset.ProblemInfo
	if problems, err = srv.problemsetService.ListProblem(problemsetId, sess.AuthUserId); err != nil {
		return
	}
	var entries []ProblemsetListProblemEntry = make([]ProblemsetListProblemEntry, len(problems))
	for i, problem := range problems {
		entries[i] = ProblemsetListProblemEntry{
			Name:  problem.Name,
			Title: problem.Title,
		}
	}
	writeResponse(w, ProblemsetListProblemResponse{Problems: entries}, sess)
}

type ProblemsetViewProblemRequest struct {
	Name string `json:"name"`
}
type ProblemsetViewProblemResponse struct {
	Statement string `json:"statement"`
}

func (srv *ApiServer) HandleProblemsetView(w http.ResponseWriter, r *http.Request) {
	var err error
	var sess *session.Session
	defer func() {
		if err != nil {
			writeError(w, r, err, sess)
		}
	}()

	if _, sess, err = srv.ensureSession(w, r); err != nil {
		return
	}
	if sess.AuthUserId == defaultUserId {
		err = NotLoggedInError
		return
	}

	vars := mux.Vars(r)
	var problemsetId uuid.UUID
	if problemsetId, err = uuid.Parse(vars["problemset_id"]); err != nil {
		return
	}
	var req ProblemsetViewProblemRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	var p1 problemset.ProblemInfo
	if p1, err = srv.problemsetService.ViewProblem(problemsetId, sess.AuthUserId, req.Name); err != nil {
		return
	}
	var p2 = new(judge.Problem)
	if err = srv.judgeService.GetProblemStatementInfo(p1.ProblemId, p2); err != nil {
		return
	}

	writeResponse(w, ProblemsetViewProblemResponse{
		Statement: p2.Statement,
	}, sess)
}

type ProblemsetSubmitRequest struct {
	ProblemName string `json:"problem_name"`
	Type        string `json:"type"`
	Traditional judge.TraditionalSubmission
}
type ProblemsetSubmitResponse struct {
	SubmissionId uuid.UUID `json:"submission_id"`
}

func (srv *ApiServer) HandleProblemsetSubmit(w http.ResponseWriter, r *http.Request) {
	var err error
	var sess *session.Session
	defer func() {
		if err != nil {
			writeError(w, r, err, sess)
		}
	}()

	if _, sess, err = srv.ensureSession(w, r); err != nil {
		return
	}
	if sess.AuthUserId == defaultUserId {
		err = NotLoggedInError
		return
	}

	vars := mux.Vars(r)
	var problemsetId uuid.UUID
	if problemsetId, err = uuid.Parse(vars["problemset_id"]); err != nil {
		return
	}
	var req ProblemsetSubmitRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	switch req.Type {
	case "traditional":
		var submissionId uuid.UUID
		if submissionId, err = srv.problemsetService.SubmitTraditional(problemsetId, sess.AuthUserId, req.ProblemName, req.Traditional); err != nil {
			return
		}
		writeResponse(w, ProblemsetSubmitResponse{
			SubmissionId: submissionId,
		}, sess)
		break
	default:
		err = BadRequestError
		return
	}
}
