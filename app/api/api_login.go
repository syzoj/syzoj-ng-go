package api

import (
	"context"
	"encoding/json"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
)

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (srv *ApiServer) Handle_Login(c *ApiContext) (apiErr ApiError) {
	var err error
	var req LoginRequest
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
		const loginCheck = `
query LoginCheck($userName: string, $password: string) {
	user(func: eq(user.username, $userName)) {
		uid
		user.username
		check: checkpwd(user.password, $password)
	}
}
`
		var apiResponse *dgo_api.Response
		if apiResponse, err = t.T.QueryWithVars(ctx, loginCheck, map[string]string{"$userName": req.UserName, "$password": req.Password}); err != nil {
			return
		}
		type response struct {
			User []*User `json:"user"`
		}
		var resp response
		if err = json.Unmarshal(apiResponse.Json, &resp); err != nil {
			return
		}
		if len(resp.User) == 0 {
			apiErr = ErrUserNotFound
			return
		}
		var user = resp.User[0]
		if !user.Check {
			apiErr = ErrPasswordIncorrect
			return
		}
		if _, err = t.T.Mutate(ctx, &dgo_api.Mutation{
			Set: []*dgo_api.NQuad{
				{
					Subject:   sess.Uid,
					Predicate: "session.auth_user",
					ObjectId:  user.Uid,
				},
			},
		}); err != nil {
			return
		}
		t.Defer(func() {
			c.sessResponse.UserName = user.UserName
			c.sessResponse.LoggedIn = true
			writeResponse(c, nil)
		})
		return
	})
	if err != nil {
		return internalServerError(err)
	}
	return
}
