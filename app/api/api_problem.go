package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/syzoj/syzoj-ng-go/app/judge"
	"github.com/syzoj/syzoj-ng-go/app/session"
)

type CreateProblemRequest struct {
	Statement judge.ProblemStatement `json:"statement"`
}
type CreateProblemResponse struct {
	ProblemId uuid.UUID `json:"problem_id"`
}

func (s *ApiServer) HandleProblemCreate(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, err)
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
		Owner:     sess.AuthUserId,
		Statement: req.Statement,
	}
	var problemId uuid.UUID
	if problemId, err = s.judgeService.CreateProblem(&info); err != nil {
		return
	}
	writeResponseWithSession(w, CreateProblemResponse{ProblemId: problemId}, sess)
}
