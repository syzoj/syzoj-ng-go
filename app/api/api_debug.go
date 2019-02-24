package api

import (
	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Debug_Contest_Submit(c *ApiContext) ApiError {
	var err error
	req := new(model.DebugContestSubmitRequest)
	if err = c.GetBody(req); err != nil {
		return badRequestError(err)
	}
	contest := c.Server().c.GetContest(model.MustGetObjectID(req.ContestId))
	userId := model.MustGetObjectID(req.UserId)
	submissionId := model.MustGetObjectID(req.SubmissionId)
	player := contest.GetPlayerById(userId)
	contest.AppendSubmission(player, req.GetName(), submissionId)
	return nil
}
