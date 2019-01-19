package api

import (
	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/valyala/fastjson"
)

func Handle_Nav_Logout(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	if _, err = c.Dgraph().NewTxn().Mutate(c.Context(), &dgo_api.Mutation{
		Del: []*dgo_api.NQuad{{
			Subject:     c.Session.SessUid,
			Predicate:   "session.auth_user",
			ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_DefaultVal{DefaultVal: "_STAR_ALL"}},
		}},
		CommitNow: true,
	}); err != nil {
		return internalServerError(err)
	}
	if err = c.SessionReload(); err != nil {
		return internalServerError(err)
	}
	c.SendValue(fastjson.NewNull())
	return
}
