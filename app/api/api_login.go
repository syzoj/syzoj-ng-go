package api

import (
	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/valyala/fastjson"
)

func Handle_Login(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	if c.Session.LoggedIn() {
		return ErrAlreadyLoggedIn
	}
	var body *fastjson.Value
	if body, err = c.GetBody(); err != nil {
		return badRequestError(err)
	}
	userName := string(body.GetStringBytes("username"))
	password := string(body.GetStringBytes("password"))
	var dgValue *fastjson.Value
	if dgValue, err = c.Query(LoginQuery, map[string]string{"$userName": userName, "$password": password}); err != nil {
		return internalServerError(err)
	}
	if len(dgValue.GetArray("user")) == 0 {
		return ErrUserNotFound
	}
	userVal := dgValue.Get("user", "0")
	if !userVal.GetBool("check") {
		return ErrPasswordIncorrect
	}
	if _, err = c.Dgraph().NewTxn().Mutate(c.Context(), &dgo_api.Mutation{
		Set: []*dgo_api.NQuad{
			{
				Subject:   c.Session.SessUid,
				Predicate: "session.auth_user",
				ObjectId:  string(userVal.GetStringBytes("uid")),
			},
		},
		CommitNow: true,
	}); err != nil {
		return internalServerError(err)
	}
	if err = c.SessionReload(); err != nil {
		return internalServerError(err)
	}
	c.SendValue(fastjson.NewNull())
	return
}
