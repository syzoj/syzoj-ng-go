package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgraph-io/dgo"
	dgo_api "github.com/dgraph-io/dgo/protos/api"
	dgo_y "github.com/dgraph-io/dgo/y"
	"github.com/google/uuid"
)

type Session struct {
	Uid       string  `json:"uid",omitempty`
	AuthUser  []*User `json:"session.auth_user",omitempty`
	UserAgent string  `json:"session.user_agent",omitempty`
	RemoteIP  string  `json:"session.remote_ip",omitempty`
}

type User struct {
	Uid      string `json:"uid",omitempty`
	UserName string `json:"user.username",omitempty`
	Check    bool   `json:"check",omitempty`
}

type Problem struct {
	Uid        string    `json:"uid",omitempty`
	Id         uuid.UUID `json:"problem.id",omitempty`
	Title      string    `json:"problem.title",omitempty`
	Statement  string    `json:"problem.statement",omitempty`
	Token      string    `json:"problem.token",omitempty`
	CreateTime time.Time `json:"problem.create_time",omitempty`
	Owner      []*User   `json:"problem.owner",omitempty`
}

type DgraphTransaction struct {
	T   *dgo.Txn
	def []func()
}

func (t *DgraphTransaction) Defer(f func()) {
	t.def = append(t.def, f)
}

func (srv *ApiServer) withDgraphTransaction(ctx context.Context, f func(context.Context, *DgraphTransaction) error) (err error) {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var abort bool
			abort = func() bool {
				t := DgraphTransaction{
					T: srv.dgraph.NewTxn(),
				}
				defer t.T.Discard(ctx)
				if err = f(ctx, &t); err != nil {
					return true
				}
				if err = t.T.Commit(ctx); err == nil {
					for _, def := range t.def {
						def()
					}
					return true
				} else if err == dgo_y.ErrAborted {
					return false
				}
				return true
			}()
			if abort {
				return
			}
			break
		}
	}
}
func (srv *ApiServer) withDgraphReadOnly(ctx context.Context, f func(context.Context, *DgraphTransaction) error) (err error) {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			var abort bool
			abort = func() bool {
				t := DgraphTransaction{
					T: srv.dgraph.NewReadOnlyTxn(),
				}
				defer t.T.Discard(ctx)
				if err = f(ctx, &t); err != nil {
					return true
				}
				if err = t.T.Commit(ctx); err == nil {
					for _, def := range t.def {
						def()
					}
					return true
				} else if err == dgo_y.ErrAborted {
					return false
				}
				return true
			}()
			if abort {
				return
			}
			break
		}
	}
}

func (srv *ApiServer) getSession(ctx context.Context, c *ApiContext, t *DgraphTransaction) (sess *Session, err error) {
	var claimedToken string
	if cookie, err := c.r.Cookie("SYZOJSESSION"); err == nil {
		claimedToken = cookie.Value
	}
	if len(claimedToken) != 32 {
		claimedToken = ""
	}
	const q = `
query Session($token: string) {
	session(func: eq(session.token, $token)) {
		uid
		session.auth_user {
			uid
			user.username
		}
	}
}
`
	type QueryResponse struct {
		Session []*Session `json:"session'`
	}
	var apiResponse *dgo_api.Response
	if apiResponse, err = t.T.QueryWithVars(ctx, q, map[string]string{"$token": claimedToken}); err != nil {
		return
	}
	var response QueryResponse
	if err = json.Unmarshal(apiResponse.Json, &response); err != nil {
		return
	}
	if len(response.Session) == 0 {
		var b [16]byte
		if _, err = rand.Read(b[:]); err != nil {
			return
		}
		var token = hex.EncodeToString(b[:])
		var apiAssigned *dgo_api.Assigned
		if apiAssigned, err = t.T.Mutate(ctx, &dgo_api.Mutation{
			Set: []*dgo_api.NQuad{
				{
					Subject:     "_:session",
					Predicate:   "session.token",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: token}},
				},
			},
		}); err != nil {
			return
		}
		sess = new(Session)
		sess.Uid = apiAssigned.Uids["session"]
		t.Defer(func() {
			http.SetCookie(c.w, &http.Cookie{
				Name:     "SYZOJSESSION",
				HttpOnly: true,
				Path:     "/",
				Value:    token,
				Expires:  time.Now().Add(time.Hour * 24 * 30),
			})
		})
	} else {
		sess = response.Session[0]
	}
	t.Defer(func() {
		c.sessResponse = new(SessionResponse)
	})
	if len(sess.AuthUser) != 0 {
		curUserName := sess.AuthUser[0].UserName
		t.Defer(func() {
			c.sessResponse.LoggedIn = true
			c.sessResponse.UserName = curUserName
		})
	}
	return
}
