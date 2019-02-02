package api

import (
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/core"
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
	var resp *core.ContestRegister1Resp
	if resp, err = c.Server().c.Action_Contest_Register(c.Context(), &core.ContestRegister1{
		UserId:    c.Session.AuthUserUid,
		ContestId: contestId,
	}); err != nil {
		switch err {
		case core.ErrAlreadyRegistered:
			return ErrAlreadyRegistered
		case core.ErrContestNotRunning:
			return ErrContestNotFound
		default:
			return internalServerError(err)
		}
	}
	_ = resp
	arena := new(fastjson.Arena)
	c.SendValue(arena.NewNull())
	return
}
