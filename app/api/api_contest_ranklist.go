package api

import (
	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Contest_Ranklist(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	contestId := model.MustDecodeObjectID(vars["contest_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	contest := c.Server().c.GetContest(contestId)
	ranklist := contest.GetRanklist()
	resp := new(model.ContestRanklistResponse)
	resp.Ranklist = ranklist
	c.SendValue(resp)
	return nil
}
