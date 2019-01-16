package api

import (
	"context"
	"encoding/json"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (srv *ApiServer) Handle_Submission_View(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := mux.Vars(c.r)
	var submissionId = uuid.MustParse(vars["submission_id"])
	err = srv.withDgraphTransaction(c.r.Context(), func(ctx context.Context, t *DgraphTransaction) (err error) {
		const SubmissionViewQuery = `
query SubmissionViewQuery($submissionId: string) {
	submission(func: eq(submission.id, $submissionId)) @normalize {
		status: submission.status
		message: submission.message
		score: submission.score
		language: submission.language
		code: submission.code
		submit_time: submission.submit_time
		submission.owner {
			submitter_name: user.username
		}
		submission.problem {
			problem_id: problem.id
			problem_title: problem.title
		}
	}
}
`
		var apiResponse *dgo_api.Response
		if apiResponse, err = t.T.QueryWithVars(ctx, SubmissionViewQuery, map[string]string{"$submissionId": submissionId.String()}); err != nil {
			return
		}
		type response struct {
			Submission []json.RawMessage `json:"submission"`
		}
		var resp response
		if err = json.Unmarshal(apiResponse.Json, &resp); err != nil {
			return
		}
		if len(resp.Submission) == 0 {
			t.Defer(func() {
				apiErr = ErrSubmissionNotFound
			})
			return
		}
		t.Defer(func() {
			writeResponse(c, resp.Submission[0])
		})
		return
	})
	if err != nil {
		return internalServerError(err)
	}
	return
}
