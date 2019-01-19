package api

import (
	"github.com/valyala/fastjson"
)

func Handle_Problems(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	var val *fastjson.Value
	if val, err = c.Query(ProblemsQuery, nil); err != nil {
		return
	}
	c.SendValue(val.Get("problems"))
	return
}
