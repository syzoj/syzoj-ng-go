package main

import (
	"encoding/json"
	"fmt"
	"time"

	httplib "github.com/syzoj/syzoj-ng-go/lib/fasthttp"
	"github.com/syzoj/syzoj-ng-go/util"
	"github.com/valyala/fasthttp"
)

type ProblemInfo struct {
	Title     string `json:"title"`
	Statement string `json:"statement"`
}

type ProblemNewRequest struct {
	Info *ProblemInfo `json:"info,omitempty"`
	Tags []string     `json:"tags,omitempty"`
}

func (app *App) postProblemNew(ctx *fasthttp.RequestCtx) {
	var req ProblemNewRequest
	if err := httplib.ReadBodyJSON(ctx, &req); err != nil {
		httplib.BadRequest(ctx, err)
		return
	}
	if req.Info == nil {
		httplib.BadRequest(ctx, fmt.Errorf("Missing info field"))
		return
	}
	sess, err := app.getSession(ctx)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	if sess == nil || sess.CurrentUser == nil {
		httplib.SendError(ctx, "Not logged in")
		return
	}
	userUid := sess.CurrentUser.UserUid
	problemUid := util.RandomString(16)
	problemInfoPayload, err := json.Marshal(req.Info)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	_, err = app.dbProblem.ExecContext(ctx, "INSERT INTO problems (uid, owner_user_uid, info) VALUES (?, ?, ?)", problemUid, userUid, problemInfoPayload)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	httplib.SendJSON(ctx, map[string]interface{}{
		"problem": map[string]interface{}{
			"uid": problemUid,
		},
	})
}

func (app *App) getProblemInfo(ctx *fasthttp.RequestCtx) {
	problemUid := ctx.UserValue("uid").(string)
	val, err := app.waitForCache(ctx, fmt.Sprintf("problem:%s:info", problemUid), time.Second*5, func() {
		app.automationCli.Trigger(map[string]interface{}{
			"tags": []string{"cache/problem/*/info/request"},
			"problem": map[string]interface{}{
				"uid": problemUid,
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
