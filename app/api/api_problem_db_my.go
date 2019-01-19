package api

import (
	"github.com/valyala/fastjson"
)

func Handle_ProblemDb_My(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	var dgValue *fastjson.Value
	if dgValue, err = c.Query(MyProblemQuery, map[string]string{"$userId": c.Session.AuthUserUid}); err != nil {
		return internalServerError(err)
	}
	c.SendValue(dgValue)
	return
}
