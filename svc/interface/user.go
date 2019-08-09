package main

import (
	"encoding/json"
	"fmt"
	"time"

	httplib "github.com/syzoj/syzoj-ng-go/lib/fasthttp"
	"github.com/valyala/fasthttp"
)

func (app *App) getUserCurrent(ctx *fasthttp.RequestCtx) {
	sess, err := app.getSession(ctx)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	if sess == nil || sess.CurrentUser == nil {
		httplib.SendError(ctx, "Not logged in")
		return
	}
	uid := sess.CurrentUser.UserUid
	val, err := app.waitForCache(ctx, fmt.Sprintf("user:%s:self", uid), time.Second*5, func() {
		app.automationCli.Trigger(map[string]interface{}{
			"tags": []string{"cache/user/*/self/request"},
			"user": map[string]interface{}{
				"uid": uid,
			},
		})
	})
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	httplib.SendJSON(ctx, map[string]interface{}{
		"data": json.RawMessage(val),
	})
}
