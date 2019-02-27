package api

import (
	"github.com/syzoj/syzoj-ng-go/app/model"

	"github.com/golang/protobuf/ptypes/empty"
)

func Handle_Contest_Register(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	contestId := model.MustDecodeObjectID(vars["contest_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	contest := c.Server().c.GetContest(contestId)
	if contest == nil {
		return ErrContestNotFound
	}
	if !contest.RegisterPlayer(c.Session.AuthUserUid) {
		return ErrGeneral
	}
	c.SendValue(new(empty.Empty))
	return nil
}
