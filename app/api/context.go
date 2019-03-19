package api

import (
	"net/http"
	"net/url"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type ApiContext struct {
	res     http.ResponseWriter
	req     *http.Request
	Session *Session
	Server  *ApiServer
}

func (c *ApiContext) Vars() map[string]string {
	return mux.Vars(c.req)
}

func (c *ApiContext) FormValue(name string) string {
	return c.req.FormValue(name)
}

func (c *ApiContext) Form() url.Values {
	c.req.ParseForm()
	return c.req.Form
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

func (c *ApiContext) getSessionVal() *model.ResponseSession {
	return &model.ResponseSession{
		UserName: proto.String(c.Session.AuthUserUserName),
		LoggedIn: proto.Bool(c.Session.LoggedIn()),
	}
}

func (c *ApiContext) SendError(err ApiError) {
	if ierr, ok := err.(internalServerErrorType); ok {
		log.Errorf("Error handling request %s: %s", c.req.URL, ierr.Err)
	} else {
		log.Infof("Failed to handle request %s: %s", c.req.URL, err)
	}
	resp := new(model.Response)
	resp.Error = proto.String(err.Error())
	if c.Session != nil {
		resp.Session = c.getSessionVal()
	}
	err2 := jsonMarshaler.Marshal(c.res, resp)
	if err2 != nil {
		log.WithField("error", err2).Warning("Failed to write error")
	}
}

func (c *ApiContext) SendValue(val proto.Message) {
	var err error
	resp := &model.Response{}
	resp.Data, err = ptypes.MarshalAny(val)
	if err != nil {
		log.WithField("error", err).Warning("Failed to write response")
		return
	}
	if c.Session != nil {
		resp.Session = c.getSessionVal()
	}
	err = jsonMarshaler.Marshal(c.res, resp)
	if err != nil {
		log.WithField("error", err).Warning("Failed to write response")
		return
	}
}

func (c *ApiContext) GetBody(msg proto.Message) error {
	return jsonUnmarshaler.Unmarshal(c.req.Body, msg)
}

func (c *ApiContext) UpgradeWebSocket() (*websocket.Conn, error) {
	return c.Server.wsUpgrader.Upgrade(c.res, c.req, nil)
}
