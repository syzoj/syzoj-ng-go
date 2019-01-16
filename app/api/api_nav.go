package api

import (
	"context"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
)

func (srv *ApiServer) Handle_Nav_Logout(c *ApiContext) (apiErr ApiError) {
	var err error
	err = srv.withDgraphTransaction(c.r.Context(), func(ctx context.Context, t *DgraphTransaction) (err error) {
		var sess *Session
		if sess, err = srv.getSession(ctx, c, t); err != nil {
			return
		}
		if len(sess.AuthUser) == 0 {
			apiErr = ErrNotLoggedIn
			return
		}
		if _, err = t.T.Mutate(ctx, &dgo_api.Mutation{
			Del: []*dgo_api.NQuad{
				{
					Subject:     sess.Uid,
					Predicate:   "session.auth_user",
					ObjectId:    "",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_DefaultVal{DefaultVal: "_STAR_ALL"}},
				},
			},
		}); err != nil {
			return
		}
		t.Defer(func() {
			c.sessResponse.LoggedIn = false
			c.sessResponse.UserName = ""
			writeResponse(c, nil)
		})
		return
	})
	if err != nil {
		return internalServerError(err)
	}
	return
}
