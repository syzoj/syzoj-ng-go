package main

import (
	"fmt"
	"regexp"

	httplib "github.com/syzoj/syzoj-ng-go/lib/fasthttp"
	"github.com/syzoj/syzoj-ng-go/util"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (app *App) postUserLogin(ctx *fasthttp.RequestCtx) {
	var req LoginRequest
	if err := httplib.ReadBodyJSON(ctx, &req); err != nil {
		httplib.BadRequest(ctx, err)
		return
	}
	sess, err := app.getSession(ctx)
	if err == nil && sess != nil && sess.CurrentUser != nil {
		httplib.SendError(ctx, "Already logged in")
		return
	}
	var uid string
	var password []byte
	err = app.dbUser.QueryRowContext(ctx, "SELECT uid, password FROM users WHERE username=?", req.UserName).Scan(&uid, &password)
	if err != nil {
		httplib.SendError(ctx, "User doesn't exist")
		return
	}
	if bcrypt.CompareHashAndPassword(password, []byte(req.Password)) != nil {
		httplib.SendError(ctx, "Password doesn't match")
		go app.automationCli.Trigger(map[string]interface{}{
			"tags":     []string{"user/*/login-failure", fmt.Sprintf("user/%s/login-failure", uid)},
			"user_uid": uid,
		})
		return
	}

	sess = &Session{}
	sess.Expire = 3600 * 24 * 30
	sess.CurrentUser = &SessionUser{
		UserUid: uid,
	}
	if err := app.newSession(ctx, sess); err != nil {
		httplib.SendError(ctx, "Internal server error")
		return
	}
	go app.automationCli.Trigger(map[string]interface{}{
		"tags":     []string{"user/*/login-success", fmt.Sprintf("user/%s/login-success", uid)},
		"user_uid": uid,
	})
}

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

var usernameRegex = regexp.MustCompile("[0-9A-Za-z_]{3,16}")

func (app *App) postUserRegister(ctx *fasthttp.RequestCtx) {
	var req RegisterRequest
	if err := httplib.ReadBodyJSON(ctx, &req); err != nil {
		httplib.BadRequest(ctx, err)
		return
	}
	sess, err := app.getSession(ctx)
	if err == nil && sess != nil && sess.CurrentUser != nil {
		httplib.SendError(ctx, "Already logged in")
		return
	}
	uid := util.RandomString(15)
	if !usernameRegex.MatchString(req.UserName) {
		httplib.SendError(ctx, "Invalid username")
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 0)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	if _, err := app.dbUser.ExecContext(ctx, "INSERT INTO users (uid, username, password) VALUES (?, ?, ?)", uid, req.UserName, passwordHash); err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	sess = &Session{}
	sess.Expire = 3600 * 24 * 30
	sess.CurrentUser = &SessionUser{
		UserUid: uid,
	}
	if err := app.newSession(ctx, sess); err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	go app.automationCli.Trigger(map[string]interface{}{
		"tags":     []string{"user/*/register", fmt.Sprintf("user/%s/register", uid)},
		"user_uid": uid,
	})
}
