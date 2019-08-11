package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	httplib "github.com/syzoj/syzoj-ng-go/lib/fasthttp"
	"github.com/syzoj/syzoj-ng-go/util"
	"github.com/valyala/fasthttp"
)

type JudgeEnqueueRequest struct {
	Info json.RawMessage `json:"info"`
}

// Task states
const stEnqueue string = "ENQUEUE"
const stJudging = "JUDGING"
const stDone = "DONE"

func (app *App) postEnqueue(ctx *fasthttp.RequestCtx) {
	var req JudgeEnqueueRequest
	if err := httplib.ReadBodyJSON(ctx, &req); err != nil {
		httplib.BadRequest(ctx, err)
		return
	}
	if req.Info == nil {
		httplib.BadRequest(ctx, fmt.Errorf("Missing info field"))
		return
	}
	taskUid := util.RandomString(16)
	conn, err := app.redisSession.GetContext(ctx)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	defer conn.Close()
	_, err = conn.Do("SET", fmt.Sprintf("judge-queue:task:%s:state", taskUid), stEnqueue, "EX", 3600*24)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	_, err = conn.Do("SET", fmt.Sprintf("judge-queue:task:%s:info", taskUid), req.Info, "EX", 3600*24)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	score := time.Now().UnixNano()
	_, err = conn.Do("ZADD", "judge-queue:queue", score, taskUid)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	httplib.SendJSON(ctx, map[string]interface{}{
		"task": map[string]interface{}{
			"uid": taskUid,
		},
	})
	go app.automationCli.Trigger(map[string]interface{}{
		"tags": []string{"judge-queue/task/*/enqueue"},
		"task": map[string]interface{}{
			"uid":  taskUid,
			"info": req.Info,
		},
	})
}

func (app *App) postFetch(ctx *fasthttp.RequestCtx) {
	timeout := ctx.QueryArgs().GetUintOrZero("timeout")
	if timeout < 0 || timeout > 60 {
		timeout = 60
	}
	conn, err := app.redisSession.GetContext(ctx)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	defer conn.Close()
	vals, err := redis.Values(conn.Do("BZPOPMIN", "judge-queue:queue", timeout))
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	taskUid, _ := redis.String(vals[1], nil)
	_, err = conn.Do("SET", fmt.Sprintf("judge-queue:task:%s:state", taskUid), stJudging)
	if err != nil {
		if err == redis.ErrNil {
			ctx.SetStatusCode(204)
		} else {
			httplib.SendInternalError(ctx, err)
		}
		return
	}
	httplib.SendJSON(ctx, map[string]interface{}{
		"task": map[string]interface{}{
			"uid": taskUid,
		},
	})
	go app.automationCli.Trigger(map[string]interface{}{
		"tags": []string{"judge-queue/task/*/judging"},
		"task": map[string]interface{}{
			"uid": taskUid,
		},
	})
}

type TaskHandleRequest struct {
	Result json.RawMessage `json:"result"`
}

func (app *App) postTaskHandle(ctx *fasthttp.RequestCtx) {
	taskUid := ctx.UserValue("uid").(string)
	var req TaskHandleRequest
	if err := httplib.ReadBodyJSON(ctx, &req); err != nil {
		httplib.BadRequest(ctx, err)
		return
	}
	if req.Result == nil {
		httplib.BadRequest(ctx, fmt.Errorf("Missing result field"))
		return
	}
	conn, err := app.redisSession.GetContext(ctx)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	defer conn.Close()
	res, err := redis.String(conn.Do("GET", fmt.Sprintf("judge-queue:task:%s:state", taskUid)))
	if err == redis.ErrNil {
		httplib.NotFound(ctx, fmt.Errorf("Task not found"))
		return
	} else if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	if res != stJudging {
		httplib.SendError(ctx, fmt.Sprintf("Expected state %s, got %s", stJudging, res))
		return
	}
	_, err = conn.Do("SET", fmt.Sprintf("judge-queue:task:%s:result", taskUid), req.Result)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	_, err = conn.Do("SET", fmt.Sprintf("judge-queue:task:%s:state", taskUid), stDone)
	if err != nil {
		httplib.SendInternalError(ctx, err)
		return
	}
	go app.automationCli.Trigger(map[string]interface{}{
		"tags": []string{"judge-queue/task/*/done"},
		"task": map[string]interface{}{
			"uid": taskUid,
		},
	})
	ctx.SetStatusCode(204)
}
