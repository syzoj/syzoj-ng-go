package api

import (
	"time"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/valyala/fastjson"
)

func Handle_Register(c *ApiContext) (apiErr ApiError) {
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
	if !checkUserName(userName) {
		return ErrInvalidUserName
	}
	password := string(body.GetStringBytes("password"))
	if err = c.DgraphTransaction(func(t *DgraphTransaction) (err error) {
		var dgResponse *dgo_api.Response
		if dgResponse, err = t.T.QueryWithVars(c.Context(), CheckUserNameQuery, map[string]string{"$userName": userName}); err != nil {
			return
		}
		dgParser := c.GetParser()
		defer c.PutParser(dgParser)
		var val *fastjson.Value
		if val, err = dgParser.ParseBytes(dgResponse.Json); err != nil {
			panic(err)
		}
		if len(val.GetArray("user")) != 0 {
			t.Defer(func() {
				apiErr = ErrDuplicateUserName
			})
			return
		}
		var timeNow []byte
		timeNow, err = time.Now().MarshalBinary()
		if err != nil {
			panic(err)
		}
		if _, err = t.T.Mutate(c.Context(), &dgo_api.Mutation{
			Set: []*dgo_api.NQuad{
				{
					Subject:     "_:user",
					Predicate:   "user.username",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: userName}},
				},
				{
					Subject:     "_:user",
					Predicate:   "user.password",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: password}},
				},
				{
					Subject:     "_:user",
					Predicate:   "user.register_time",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_DatetimeVal{DatetimeVal: timeNow}},
				},
			},
		}); err != nil {
			return
		}
		log.WithField("username", userName).Info("Created account")
		t.Defer(func() {
			c.SendValue(fastjson.NewNull())
		})
		return
	}); err != nil {
		return internalServerError(err)
	}
	return
}
