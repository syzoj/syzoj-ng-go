package api

import (
	"context"
	"encoding/json"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
)

type MySubmissionResponse struct {
	Submissions json.RawMessage `json:"submissions"`
}

func (srv *ApiServer) Handle_Submission_My(c *ApiContext) (apiErr ApiError) {
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
		const MySubmissionQuery = `
query MySubmissionQuery($userId: string) {
	submissions(func: uid($userId)) @normalize {
		~submission.owner {
			submission_id: submission.id
			submission_status: submission.status
			submit_time: submission.submit_time
			submission.problem {
				problem_id: problem.id
				problem_title: problem.title
			}
		}
	}
}
`
		var apiResponse *dgo_api.Response
		if apiResponse, err = t.T.QueryWithVars(ctx, MySubmissionQuery, map[string]string{"$userId": sess.AuthUser[0].Uid}); err != nil {
			return
		}
		type response struct {
			Submissions json.RawMessage `json:"submissions"`
		}
		var resp response
		if err = json.Unmarshal(apiResponse.Json, &resp); err != nil {
			return
		}
		var myresp = MySubmissionResponse{Submissions: resp.Submissions}
		t.Defer(func() {
			writeResponse(c, myresp)
		})
		return
	})
	if err != nil {
		return internalServerError(err)
	}
	return
}
