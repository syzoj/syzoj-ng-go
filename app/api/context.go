package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/valyala/fastjson"
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

func (c *ApiContext) Server() *ApiServer {
	return c.srv
}

func (c *ApiContext) FormValue(name string) string {
	return c.req.FormValue(name)
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

func (c *ApiContext) getSessionVal(arena *fastjson.Arena) *fastjson.Value {
	val := arena.NewObject()
	val.Set("user_name", arena.NewString(c.Session.AuthUserUserName))
	if c.Session.LoggedIn() {
		val.Set("logged_in", arena.NewTrue())
	} else {
		val.Set("logged_in", arena.NewFalse())
	}
	//godoc.org/github.com/mongodb/mongo-go-driver/mongo#Databasehttps://godoc.org/github.com/mongodb/mongo-go-driver/mongo#Databasehttps://godoc.org/github.com/mongodb/mongo-go-driver/mongo#Databasegithub.com/dgraph-io/dgo"
	return val
}

func (c *ApiContext) SendError(err ApiError) {
	if ierr, ok := err.(internalServerErrorType); ok {
		log.Errorf("Error handling request %s: %s", c.req.URL, ierr.Err)
	} else {
		log.Infof("Failed to handle request %s: %s", c.req.URL, err)
	}
	arena := new(fastjson.Arena)
	val := arena.NewObject()
	val.Set("error", arena.NewString(err.Error()))
	if c.Session != nil {
		val.Set("session", c.getSessionVal(arena))
	}
	_, err2 := c.res.Write(val.MarshalTo(nil))
	if err2 != nil {
		log.WithField("error", err2).Warning("Failed to write error")
	}
}

func (c *ApiContext) SendValue(val *fastjson.Value) {
	arena := new(fastjson.Arena)
	mval := arena.NewObject()
	mval.Set("data", val)
	if c.Session != nil {
		mval.Set("session", c.getSessionVal(arena))
	}
	data := mval.MarshalTo(nil)
	c.SetHeader("Content-Length", strconv.Itoa(len(data)))
	_, err := c.res.Write(data)
	if err != nil {
		log.WithField("error", err).Warning("Failed to write response")
	}
}

func (c *ApiContext) Context() context.Context {
	return c.req.Context()
}

func (c *ApiContext) GetBody() (*fastjson.Value, error) {
	var err error
	buf := bytes.Buffer{}
	if _, err = io.Copy(&buf, c.req.Body); err != nil {
		return nil, err
	}
	return fastjson.ParseBytes(buf.Bytes())
}

func (c *ApiContext) UpgradeWebSocket() (*websocket.Conn, error) {
	return c.srv.wsUpgrader.Upgrade(c.res, c.req, nil)
}
