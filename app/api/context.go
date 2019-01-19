package api

import (
	"bytes"
	"context"
	"io"
	"net/http"

	"github.com/dgraph-io/dgo"
	dgo_api "github.com/dgraph-io/dgo/protos/api"
	dgo_y "github.com/dgraph-io/dgo/y"
	"github.com/gorilla/mux"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/judge"
)

type ApiContext struct {
	res     http.ResponseWriter
	req     *http.Request
	Session *Session
	srv     *ApiServer
}

func (c *ApiContext) Vars() map[string]string {
	return mux.Vars(c.req)
}

func (c *ApiContext) JudgeService() judge.Service {
	return c.srv.judgeService
}

func (c *ApiContext) GetCookie(name string) string {
	cookie, err := c.req.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func (c *ApiContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.res, cookie)
}

func (c *ApiContext) GetHeader(name string) string {
	return c.req.Header.Get(name)
}

func (c *ApiContext) SetHeader(name string, value string) {
	c.res.Header().Add(name, value)
}

func (c *ApiContext) SendError(err ApiError) {
	if ierr, ok := err.(internalServerErrorType); ok {
		log.Errorf("Error handling request %s: %s", c.req.URL, ierr.Err)
	} else {
		log.Infof("Failed to handle request %s: %s", c.req.URL, err)
	}
	parser := c.GetParser()
	defer c.PutParser(parser)
	val := parser.NewObject(map[string]*fastjson.Value{
		"error": parser.NewString(err.Error()),
	})
	if c.Session != nil {
		val.Set("session", parser.NewObject(map[string]*fastjson.Value{
			"user_name": parser.NewString(c.Session.AuthUserUserName),
			"logged_in": parser.NewBool(c.Session.LoggedIn()),
		}))
	}
	_, err2 := c.res.Write(val.MarshalTo(nil))
	if err2 != nil {
		log.WithField("error", err2).Warning("Failed to write error")
	}
}

func (c *ApiContext) SendValue(val *fastjson.Value) {
	parser := c.GetParser()
	defer c.PutParser(parser)
	mval := parser.NewObject(map[string]*fastjson.Value{
		"data": val,
		"session": parser.NewObject(map[string]*fastjson.Value{
			"user_name": parser.NewString(c.Session.AuthUserUserName),
			"logged_in": parser.NewBool(c.Session.LoggedIn()),
		}),
	})
	_, err := c.res.Write(mval.MarshalTo(nil))
	if err != nil {
		log.WithField("error", err).Warning("Failed to write response")
	}
}

func (c *ApiContext) Dgraph() *dgo.Dgraph {
	return c.srv.dgraph
}

func (c *ApiContext) Context() context.Context {
	return c.req.Context()
}

func (c *ApiContext) GetParser() *fastjson.Parser {
	return c.srv.parserPool.Get()
}

func (c *ApiContext) PutParser(p *fastjson.Parser) {
	c.srv.parserPool.Put(p)
}

func (c *ApiContext) GetBody() (*fastjson.Value, error) {
	var err error
	buf := bytes.Buffer{}
	if _, err = io.Copy(&buf, c.req.Body); err != nil {
		return nil, err
	}
	return fastjson.ParseBytes(buf.Bytes())
}

type DgraphTransaction struct {
	T   *dgo.Txn
	def []func()
}

func (t *DgraphTransaction) Defer(f func()) {
	t.def = append(t.def, f)
}

func (c *ApiContext) DgraphTransaction(f func(t *DgraphTransaction) error) error {
	for {
		select {
		case <-c.Context().Done():
			return c.Context().Err()
		default:
			t := DgraphTransaction{
				T: c.srv.dgraph.NewTxn(),
			}
			var err error
			err = f(&t)
			if err == dgo_y.ErrAborted {
				t.T.Discard(c.Context())
				continue
			} else if err != nil {
				t.T.Discard(c.Context())
				return err
			}
			if err = t.T.Commit(c.Context()); err == dgo_y.ErrAborted {
				t.T.Discard(c.Context())
				continue
			} else if err == nil {
				for _, def := range t.def {
					def()
				}
				return nil
			} else {
				t.T.Discard(c.Context())
				return err
			}
		}
	}
}

func (c *ApiContext) Query(q string, p map[string]string) (v *fastjson.Value, err error) {
	var dgResponse *dgo_api.Response
	dgResponse, err = c.srv.dgraph.NewReadOnlyTxn().QueryWithVars(c.Context(), q, p)
	if err != nil {
		return
	}
	return fastjson.ParseBytes(dgResponse.Json)
}
