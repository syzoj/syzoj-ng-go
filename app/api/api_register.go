package api

import (
	"context"
	"encoding/json"
	"time"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
)

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (srv *ApiServer) Handle_Register(c *ApiContext) (apiErr ApiError) {
	var err error
	var req RegisterRequest
	if err = json.NewDecoder(c.r.Body).Decode(&req); err != nil {
		return badRequestError(err)
	}
	err = srv.withDgraphTransaction(c.r.Context(), func(ctx context.Context, t *DgraphTransaction) (err error) {
		var sess *Session
		if sess, err = srv.getSession(ctx, c, t); err != nil {
			return
		}
		if len(sess.AuthUser) != 0 {
			apiErr = ErrAlreadyLoggedIn
			return
		}
		const registerCheck = `
query RegisterCheck($userName: string) {
	user(func: eq(user.username, $userName)) {
		uid
	}
}
`
		var apiResponse *dgo_api.Response
		if apiResponse, err = t.T.QueryWithVars(ctx, registerCheck, map[string]string{"$userName": req.UserName}); err != nil {
			return
		}
		type response struct {
			User []*User `json:"user"`
		}
		var resp response
		if err = json.Unmarshal(apiResponse.Json, &resp); err != nil {
			return
		}
		if len(resp.User) != 0 {
			t.Defer(func() {
				apiErr = ErrDuplicateUserName
			})
			return
		}
		var timeNow []byte
		timeNow, err = time.Now().MarshalBinary()
		if err != nil {
			return
		}
		if _, err = t.T.Mutate(ctx, &dgo_api.Mutation{
			Set: []*dgo_api.NQuad{
				{
					Subject:     "_:user",
					Predicate:   "user.username",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: req.UserName}},
				},
				{
					Subject:     "_:user",
					Predicate:   "user.password",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: req.Password}},
				},
				{
					Subject:     "_:user",
					Predicate:   "user.register_time",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_DatetimeVal{DatetimeVal: timeNow}},
				},
			},
		}); err != nil {
			return err
		}
		log.WithField("username", req.UserName).Info("Created account")
		t.Defer(func() {
			writeResponse(c, nil)
		})
		return
	})
	if err != nil {
		return internalServerError(err)
	}
	return
}
