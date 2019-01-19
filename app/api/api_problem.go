package api

import (
	"github.com/valyala/fastjson"
)

func Handle_Problem(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	problemName := vars["name"]
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	var dgVal *fastjson.Value
	if dgVal, err = c.Query(ViewProblemQuery, map[string]string{"$problemName": problemName}); err != nil {
		return internalServerError(err)
	}
	if len(dgVal.GetArray("problem")) == 0 {
		return ErrProblemNotFound
	}
	c.SendValue(dgVal.Get("problem", "0"))
	return
}
