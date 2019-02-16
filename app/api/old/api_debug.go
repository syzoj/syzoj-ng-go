package api

import (
	"github.com/valyala/fastjson"
	"io"
	"os"
)

func Handle_Debug_Submission_Enqueue(c *ApiContext) ApiError {
	vars := c.Vars()
	submissionId := DecodeObjectID(vars["submission_id"])
	c.Server().c.EnqueueSubmission(submissionId)
	arena := new(fastjson.Arena)
	c.SendValue(arena.NewNull())
	return nil
}

func Handle_Debug_Upload(c *ApiContext) ApiError {
	stream := c.Server().GetStream("token1")
	if stream == nil {
		return ErrGeneral
	}
	go func() {
		<-c.Context().Done()
		stream.Close()
	}()
	file, err := os.Create("/tmp/1.zip")
	if err != nil {
		panic(err)
	}
	n, err := io.Copy(file, stream)
	if err != nil {
		log.Error(err)
	}
	arena := new(fastjson.Arena)
	c.SendValue(arena.NewNumberInt(int(n)))
	return nil
}
