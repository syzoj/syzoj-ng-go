package api

import (
	"encoding/json"
	"context"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	dgo_api "github.com/dgraph-io/dgo/protos/api"
)

type CreateProblemRequest struct {
	Title string `json:"title"`
}
type CreateProblemResponse struct {
	ProblemId uuid.UUID `json:"problem_id"`
}

func (srv *ApiServer) HandleProblemCreate(c *ApiContext) (apiErr ApiError) {
	var err error
	var req CreateProblemRequest
	if err = json.NewDecoder(c.r.Body).Decode(&req); err != nil {
		return badRequestError(err)
	}
	err = srv.withDgraphTransaction(c.r.Context(), func(ctx context.Context, t *DgraphTransaction) (err error) {
		var sess *Session
		if sess, err = srv.getSession(ctx, c, t); err != nil {
			return
		}
		if len(sess.AuthUser) == 0 {
			apiErr = ErrNotLoggedIn
			return
		}
		var problemId uuid.UUID
		if problemId, err = uuid.NewRandom(); err != nil {
			return
		}
		if _, err = t.T.Mutate(ctx, &dgo_api.Mutation{
			Set: []*dgo_api.NQuad{
				&dgo_api.NQuad{
					Subject: "_:problem",
					Predicate: "problem.id",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: problemId.String()}},
				},
				&dgo_api.NQuad{
					Subject: "_:problem",
					Predicate: "problem.owner",
					ObjectId: sess.AuthUser[0].Uid,
				},
				&dgo_api.NQuad{
					Subject: "_:problem",
					Predicate: "problem.title",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: req.Title}},
				},
			},
		}); err != nil {
			return
		}
		t.Defer(func() {
			writeResponse(c, CreateProblemResponse{ProblemId: problemId})
		})
		return
	})
	if err != nil {
		return internalServerError(err)
	}
	return
}

type ViewProblemResponse struct {
	Title     string `json:"title"`
	Statement string `json:"statement"`
	IsOwner   bool   `json:"is_owner"`
	Token     string `json:"token"`
}

func (srv *ApiServer) HandleProblemView(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := mux.Vars(c.r)
	var problemId = uuid.MustParse(vars["problem_id"])
	err = srv.withDgraphTransaction(c.r.Context(), func(ctx context.Context, t *DgraphTransaction) (err error) {
		var sess *Session
		if sess, err = srv.getSession(ctx, c, t); err != nil {
			return
		}
		var apiResponse *dgo_api.Response
		const problemViewQuery = `
query ProblemViewQuery($problemId: string) {
	problem(func: eq(problem.id, $problemId)) {
		uid
		problem.title: problem.title@.
		problem.statement: problem.statement@.
		problem.token
		problem.owner {
			uid
		}
	}
}
`
		apiResponse, err = t.T.QueryWithVars(ctx, problemViewQuery, map[string]string{"$problemId": problemId.String()})
		if err != nil {
			return
		}
		type response struct {
			Problem []*Problem `json:"problem"`
		}
		var resp response
		if err = json.Unmarshal(apiResponse.Json, &resp); err != nil {
			return
		}
		if len(resp.Problem) == 0 {
			return ErrProblemNotFound
		}
		var problem = resp.Problem[0]
		t.Defer(func() {
			var myresp ViewProblemResponse
			myresp.Title = problem.Title
			myresp.Statement = problem.Statement
			myresp.IsOwner = len(problem.Owner) > 0 && len(sess.AuthUser) > 0 && problem.Owner[0].Uid == sess.AuthUser[0].Uid
			if myresp.IsOwner {
				myresp.Token = problem.Token
			}
			writeResponse(c, &myresp)
		})
		return
	})
	if err != nil {
		return internalServerError(err)
	}
	return
}

type ResetProblemTokenResponse struct {
	Token string `json:"token"`
}

/*
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
*/

/*
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
*/

type ProblemChangeTitleRequest struct {
	Title string `json:"title"`
}

/*
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
*/
