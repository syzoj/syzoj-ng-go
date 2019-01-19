package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/valyala/fastjson"
)

type User struct {
	Uid      string `json:"uid",omitempty`
	UserName string `json:"user.username",omitempty"`
}

type Session struct {
	SessUid          string
	AuthUserUid      string
	AuthUserUserName string
}

func (s *Session) LoggedIn() bool {
	return s.AuthUserUid != ""
}

func (c *ApiContext) SessionStart() error {
	var claimedToken = c.GetCookie("SYZOJSESSION")
	if len(claimedToken) != 32 {
		claimedToken = ""
	}

	return c.DgraphTransaction(func(t *DgraphTransaction) (err error) {
		var dgResponse *dgo_api.Response
		if dgResponse, err = t.T.QueryWithVars(c.Context(), sessionQuery, map[string]string{"$token": claimedToken}); err != nil {
			return
		}
		parser := c.GetParser()
		defer c.PutParser(parser)
		var val *fastjson.Value
		val, err = parser.ParseBytes(dgResponse.Json)
		if err != nil {
			panic(err)
		}
		var sess = new(Session)
		if len(val.GetArray("session")) == 0 {
			var b [16]byte
			if _, err = rand.Read(b[:]); err != nil {
				return
			}
			var newToken = hex.EncodeToString(b[:])
			var dgAssigned *dgo_api.Assigned
			if dgAssigned, err = t.T.Mutate(c.Context(), &dgo_api.Mutation{
				Set: []*dgo_api.NQuad{
					{
						Subject:     "_:session",
						Predicate:   "session.token",
						ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: newToken}},
					},
				},
			}); err != nil {
				return
			}
			sess.SessUid = dgAssigned.Uids["session"]
			t.Defer(func() {
				c.SetCookie(&http.Cookie{
					Name:     "SYZOJSESSION",
					HttpOnly: true,
					Path:     "/",
					Value:    newToken,
					Expires:  time.Now().Add(time.Hour * 24 * 30),
				})
			})
		} else {
			sessVal := val.Get("session", "0")
			sess.SessUid = string(sessVal.GetStringBytes("uid"))
			sess.AuthUserUid = string(sessVal.GetStringBytes("session.auth_user", "0", "uid"))
			sess.AuthUserUserName = string(sessVal.GetStringBytes("session.auth_user", "0", "user.username"))
		}
		t.Defer(func() {
			c.Session = sess
		})
		return
	})
}

func (c *ApiContext) SessionReload() (err error) {
	var dgResponse *dgo_api.Response
	if dgResponse, err = c.srv.dgraph.NewReadOnlyTxn().QueryWithVars(c.Context(), sessionByUidQuery, map[string]string{"$sessUid": c.Session.SessUid}); err != nil {
		return err
	}
	parser := c.GetParser()
	defer c.PutParser(parser)
	var val *fastjson.Value
	val, err = parser.ParseBytes(dgResponse.Json)
	if err != nil {
		panic(err)
	}
	var sess = new(Session)
	sessVal := val.Get("session", "0")
	sess.SessUid = string(sessVal.GetStringBytes("uid"))
	sess.AuthUserUid = string(sessVal.GetStringBytes("session.auth_user", "0", "uid"))
	sess.AuthUserUserName = string(sessVal.GetStringBytes("session.auth_user", "0", "user.username"))
	c.Session = sess
	return
}
