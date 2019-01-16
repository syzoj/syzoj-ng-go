package api

import (
	"context"
	"encoding/json"
	"time"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/google/uuid"
)

type CreateProblemRequest struct {
	Title string `json:"title"`
}
type CreateProblemResponse struct {
	ProblemId uuid.UUID `json:"problem_id"`
}

func (srv *ApiServer) Handle_ProblemDb_New(c *ApiContext) (apiErr ApiError) {
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
			t.Defer(func() {
				apiErr = ErrNotLoggedIn
			})
			return
		}
		var problemId uuid.UUID
		if problemId, err = uuid.NewRandom(); err != nil {
			return
		}
		var datetimeVal []byte
		if datetimeVal, err = time.Now().MarshalBinary(); err != nil {
			return
		}
		if _, err = t.T.Mutate(ctx, &dgo_api.Mutation{
			Set: []*dgo_api.NQuad{
				{
					Subject:     "_:problem",
					Predicate:   "problem.id",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: problemId.String()}},
				},
				{
					Subject:   "_:problem",
					Predicate: "problem.owner",
					ObjectId:  sess.AuthUser[0].Uid,
				},
				{
					Subject:     "_:problem",
					Predicate:   "problem.title",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: req.Title}},
				},
				{
					Subject:     "_:problem",
					Predicate:   "problem.create_time",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_DatetimeVal{DatetimeVal: datetimeVal}},
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
