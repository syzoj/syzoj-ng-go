package api

import (
	"context"
	"encoding/json"
	"time"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ViewProblemResponse struct {
	Title     string `json:"title"`
	Statement string `json:"statement"`
	IsOwner   bool   `json:"is_owner"`
	Token     string `json:"token"`
	CanSubmit bool   `json:"can_submit"`
}

func (srv *ApiServer) Handle_ProblemDb_View(c *ApiContext) (apiErr ApiError) {
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
			t.Defer(func() {
				apiErr = ErrProblemNotFound
			})
			return
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
			myresp.CanSubmit = len(sess.AuthUser) > 0
			writeResponse(c, &myresp)
		})
		return
	})
	if err != nil {
		return internalServerError(err)
	}
	return
}

type SubmitProblemRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}
type SubmitProblemResponse struct {
	SubmissionId uuid.UUID `json:"submission_id"`
}

func (srv *ApiServer) Handle_ProblemDb_View_Submit(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := mux.Vars(c.r)
	problemId := uuid.MustParse(vars["problem_id"])
	var req SubmitProblemRequest
	if err = json.NewDecoder(c.r.Body).Decode(&req); err != nil {
		return badRequestError(err)
	}
	err = srv.withDgraphTransaction(c.r.Context(), func(ctx context.Context, t *DgraphTransaction) (err error) {
		var sess *Session
		if sess, err = srv.getSession(ctx, c, t); err != nil {
			return
		}
		if len(sess.AuthUser) == 0 {
			t.Defer(func() {
				apiErr = ErrNotLoggedIn
			})
			return
		}
		const CanSubmitQuery = `
query CanSubmitQuery($problemId: string) {
	problem(func: eq(problem.id, $problemId)) {
		uid
	}
}
`
		var apiResponse *dgo_api.Response
		apiResponse, err = t.T.QueryWithVars(ctx, CanSubmitQuery, map[string]string{"$problemId": problemId.String()})
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
			t.Defer(func() {
				apiErr = ErrProblemNotFound
			})
			return
		}
		var submissionId uuid.UUID
		if submissionId, err = uuid.NewRandom(); err != nil {
			return
		}
		var datetimeVal []byte
		if datetimeVal, err = time.Now().MarshalBinary(); err != nil {
			return
		}
		var assigned *dgo_api.Assigned
		if assigned, err = t.T.Mutate(ctx, &dgo_api.Mutation{
			Set: []*dgo_api.NQuad{
				{
					Subject:     "_:submission",
					Predicate:   "submission.id",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: submissionId.String()}},
				},
				{
					Subject:   "_:submission",
					Predicate: "submission.problem",
					ObjectId:  resp.Problem[0].Uid,
				},
				{
					Subject:     "_:submission",
					Predicate:   "submission.language",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: req.Language}},
				},
				{
					Subject:     "_:submission",
					Predicate:   "submission.code",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: req.Code}},
				},
				{
					Subject:   "_:submission",
					Predicate: "submission.owner",
					ObjectId:  sess.AuthUser[0].Uid,
				},
				{
					Subject:     "_:submission",
					Predicate:   "submission.submit_time",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_DatetimeVal{DatetimeVal: datetimeVal}},
				},
				{
					Subject:     "_:submission",
					Predicate:   "submission.status",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: "Waiting"}},
				},
			},
		}); err != nil {
			return
		}
		t.Defer(func() {
			writeResponse(c, &SubmitProblemResponse{
				SubmissionId: submissionId,
			})
		})
		t.Defer(func() {
			if err := srv.judgeService.NotifySubmission(ctx, assigned.Uids["submission"]); err != nil {
				log.Error(err)
			}
		})
		return
	})
	if err != nil {
		return internalServerError(err)
	}
	return
}
