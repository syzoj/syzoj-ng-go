package api

import (
	"github.com/valyala/fastjson"
)

func Handle_Submission_My(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	var dgVal *fastjson.Value
	if dgVal, err = c.Query(MySubmissionQuery, map[string]string{"$userId": c.Session.AuthUserUid}); err != nil {
		return
	}
	c.SendValue(fastjson.NewObject(map[string]*fastjson.Value{
		"submissions": dgVal.Get("submissions"),
	}))
	return
}
