package api

import (
	"github.com/valyala/fastjson"
	"github.com/google/uuid"
)

func Handle_Submission_View(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	var submissionId = uuid.MustParse(vars["submission_id"])
	var dgVal *fastjson.Value
	if dgVal, err = c.Query(SubmissionViewQuery, map[string]string{"$submissionId": submissionId.String()}); err != nil {
		return internalServerError(err)
	}
	if len(dgVal.GetArray("submission")) == 0 {
		return ErrSubmissionNotFound
	}
	c.SendValue(dgVal.Get("submission", "0"))
	return
}
