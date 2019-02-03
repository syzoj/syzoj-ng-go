package api

import (
	"github.com/valyala/fastjson"
)

func Handle_Contest_Register(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	contestId := DecodeObjectID(vars["contest_id"])
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	contest := c.Server().c.GetContestW(contestId)
	if contest == nil {
		return ErrContestNotFound
	}
	defer contest.Unlock()
	if !contest.RegisterPlayer(c.Session.AuthUserUid) {
		return ErrGeneral
	}
	arena := new(fastjson.Arena)
	c.SendValue(arena.NewNull())
	return
}
