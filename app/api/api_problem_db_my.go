package api

import (
	"context"
	"encoding/json"
	"time"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/google/uuid"
)

type MyProblemResponse struct {
	Problems []MyProblemSummary `json:"problems"`
}
type MyProblemSummary struct {
	Title      string    `json:"title"`
	Id         uuid.UUID `json:"id"`
	CreateTime time.Time `json:"create_time"`
}

func (srv *ApiServer) Handle_ProblemDb_My(c *ApiContext) (apiErr ApiError) {
	var err error
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
		const MyProblemQuery = `
query MyProblem($userId: string) {
	problem(func: uid($userId)) @normalize {
		~problem.owner {
			problem.id: problem.id
			problem.title: problem.title
			problem.create_time: problem.create_time
		}
	}
}
`
		var apiResponse *dgo_api.Response
		if apiResponse, err = t.T.QueryWithVars(ctx, MyProblemQuery, map[string]string{"$userId": sess.AuthUser[0].Uid}); err != nil {
			return
		}
		type response struct {
			Problem []*Problem `json:"problem"`
		}
		var resp response
		if err = json.Unmarshal(apiResponse.Json, &resp); err != nil {
			return
		}
		var myresp MyProblemResponse
		myresp.Problems = make([]MyProblemSummary, len(resp.Problem))
		for k, v := range resp.Problem {
			myresp.Problems[k] = MyProblemSummary{Title: v.Title, Id: v.Id, CreateTime: v.CreateTime}
		}
		t.Defer(func() {
			writeResponse(c, &myresp)
		})
		return
	})
	if err != nil {
		return internalServerError(err)
	}
	return
}
