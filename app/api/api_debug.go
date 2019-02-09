package api

import (
	"github.com/valyala/fastjson"
)

func Handle_Debug_Submission_Enqueue(c *ApiContext) ApiError {
	vars := c.Vars()
	submissionId := DecodeObjectID(vars["submission_id"])
	c.Server().c.EnqueueSubmission(submissionId)
	arena := new(fastjson.Arena)
	c.SendValue(arena.NewNull())
	return nil
}
