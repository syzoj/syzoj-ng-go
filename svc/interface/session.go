package main

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/syzoj/syzoj-ng-go/util"
	"github.com/valyala/fasthttp"
)

type Session struct {
	CurrentUser *SessionUser `json:"cur_user,omitempty"`
	Expire      int64        `json:"expire"` // In seconds
}

type SessionUser struct {
	UserUid string `json:"user_uid"`
}

func (app *App) getSession(ctx *fasthttp.RequestCtx) (*Session, error) {
	key := string(ctx.Request.Header.Cookie("SESSION"))
	conn, err := app.redisSession.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", fmt.Sprintf("session:%s", key)))
	if err == redis.ErrNil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	sess := &Session{}
	if err := json.Unmarshal(reply, sess); err != nil {
		log.WithError(err).Errorf("Failed to decode session %s", key)
		return nil, err
	}
	if sess.Expire != 0 {
		if _, err := conn.Do("EXPIRE", key, sess.Expire); err != nil {
			log.WithError(err).Error("Redis failure")
		}
	}
	return sess, nil
}

func (app *App) newSession(ctx *fasthttp.RequestCtx, sess *Session) error {
	key := util.RandomString(16)
	conn, err := app.redisSession.GetContext(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(sess)
	if err != nil {
		log.WithError(err).Error("Failed to encode session")
		return err
	}
	_, err = conn.Do("SET", fmt.Sprintf("session:%s", key), data)
	if err != nil {
		log.WithError(err).Error("Redis failure")
		return err
	}
	if sess.Expire != 0 {
		if _, err := conn.Do("EXPIRE", key, sess.Expire); err != nil {
			log.WithError(err).Error("Redis failure")
		}
	}
	cookie := fasthttp.AcquireCookie()
	cookie.SetKey("SESSION")
	cookie.SetValue(key)
	cookie.SetHTTPOnly(true)
	cookie.SetPath("/")
	ctx.Response.Header.SetCookie(cookie)
	fasthttp.ReleaseCookie(cookie)
	return nil
}
