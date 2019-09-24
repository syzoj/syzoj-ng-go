package app

import (
	"context"
	"encoding/json"
	"regexp"
	"sync"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/syzoj/syzoj-ng-go/lib/rediskey"
	"github.com/syzoj/syzoj-ng-go/svc/judge"
	"github.com/volatiletech/null"
)

// # Judge core
// This part of code is designed to be stateless and all state is stored in Redis for persistence.
// List of redis keys used:
// 1. rediskey.CORE_QUEUE: A Redis 5 stream where each entry is a submission. Every entry has the following keys:
//    - sid: Submission id.
//    - data: additional related data in JSON format.
//    It has a consumer group called "judger". Judgers are identified by their username in HTTP basic auth. The username is used as the Redis consumer name.
//    Redis mechanisms should be used to handle failures.
// 2. rediskey.CORE_SUBMISSION_PROGRESS: A Redis 5 stream containing submission progress. Every entry has the following keys:
//    - type: One of "reset", "progress", "done".
//    - data: JSON data if type is "progress".
// ## Judger protocol
// The judger endpoints use HTTP basic authentication. Username must match [0-9A-Za-z_]{3,16} and password is the judger token in App.JudgeToken.
// The judger connects to /judger/wait-for-task and upgrades the connection to a WebSocket. Then it receives tasks, process the task, send progress and result, and receive another task and loop.
// While processing a task, the judger sends messages with type="progress" to the broker. When the judger has done with the task, it sends a message with type="finish".
// The broker saves the result to the database and emits a type="done" event.

const GIN_JUDGER = "github.com/syzoj/syzoj-ng-go/svc/app.judger"

func (a *App) ensureQueue(ctx context.Context, queueName string) error {
	_, err := a.Redis.DoContext(ctx, "XGROUP", "CREATE", rediskey.CORE_QUEUE.Format(queueName), "judger", "$", "MKSTREAM")
	if rerr, ok := err.(redis.Error); ok {
		if rerr[:10] == "BUSYGROUP " {
			return nil
		}
	}
	return err
}

var judgerRegexp = regexp.MustCompile("[0-9A-Za-z_]{3,16}")

func (a *App) useCheckJudgeToken(c *gin.Context) {
	judger := c.Query("user")
	token := c.Query("token")
	if a.JudgeToken == "" || token != a.JudgeToken || !judgerRegexp.Match([]byte(judger)) {
		c.JSON(403, gin.H{"error": "Token doesn't match or invalid judger name (judger name must be [0-9A-Za-z_]{3,16})"})
		return
	}
	c.Set(GIN_JUDGER, judger)
}

var wsUpgrader = websocket.Upgrader{} // TODO: configure ReadBufferSize and WriteBufferSize (seems to default to 4096)

type TaskResult struct {
	Id   string          `json:"id"`
	Sid  string          `json:"sid"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// TODO: move logic to judge
func (a *App) getJudgeWaitForTask(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.WithError(err).Error("failed to upgrade websocket")
		return
	}
	defer conn.Close()

	key := rediskey.CORE_QUEUE.Format("default")
	judger := c.GetString(GIN_JUDGER)
	sema := make(chan struct{}, 1) // queue capacity is 1
	sema <- struct{}{}
	var wg sync.WaitGroup
	wg.Add(2)
	// Fetch tasks
	go func() {
		defer wg.Done()
		msgCh, errCh := a.Redis.ReadStreamGroup(ctx, key, "judger", judger, sema)
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgCh:
				sid := msg.Data["sid"]
				data, err := redis.Bytes(a.Redis.DoContext(ctx, "GET", rediskey.CORE_SUBMISSION_DATA.Format(sid)))
				if err != nil {
					log.WithError(err).Error("failed to call GET")
					return
				}
				if _, err := a.Redis.DoContext(ctx, "XADD", rediskey.CORE_SUBMISSION_PROGRESS.Format(sid), "*", "type", "reset"); err != nil {
					log.WithError(err).Error("failed to call XADD")
					return
				}
				if err := conn.WriteJSON(gin.H{
					"id":   msg.ID,
					"sid":  sid,
					"data": json.RawMessage(data),
				}); err != nil {
					log.WithError(err).Error("failed to send data")
					return
				}
			case err := <-errCh:
				if err != nil {
					log.WithError(err).Error("failed to read from Redis")
				}
				return
			}
		}
	}()

	// Receive results
	go func() {
		defer cancel()
		defer wg.Done()
		for {
			data := &TaskResult{}
			if err := conn.ReadJSON(data); err != nil {
				log.WithError(err).Error("failed to read data")
				return
			}
			skey := rediskey.CORE_SUBMISSION_PROGRESS.Format(data.Sid)
			switch data.Type {
			case "progress":
				if _, err := a.Redis.DoContext(ctx, "XADD", skey, "*", "type", "progress", "data", []byte(data.Data)); err != nil {
					log.WithError(err).Error("failed to call XADD")
					return
				}
			case "finish":
				res := &judge.Judge{}
				if err := json.Unmarshal(data.Data, res); err != nil {
					log.WithField("sid", data.Sid).WithError(err).Error("failed to parse data")
					return
				}
				if res.Type == null.IntFrom(4) {
					if err := a.JudgeService.SaveTask(ctx, data.Sid, res); err != nil {
						log.WithField("sid", data.Sid).WithError(err).Error("failed to save submission")
						return
					}
					if _, err := a.Redis.DoContext(ctx, "XADD", skey, "*", "type", "done"); err != nil {
						log.WithError(err).Error("failed to call XADD")
						return
					}
					// Keep submission progress for 24 hours
					if _, err := a.Redis.DoContext(ctx, "EXPIRE", skey, 60 * 60 * 24); err != nil {
						log.WithError(err).Error("failed to call EXPIRE")
						return
					}
					n, err := redis.Int64(a.Redis.DoContext(ctx, "XACK", key, "judger", data.Id))
					if err != nil {
						log.WithError(err).Error("failed to call XACK")
						return
					}
					for ; n > 0; n-- {
						select {
						case <-ctx.Done():
							return
						case sema <- struct{}{}:
						}
					}
				}
			}
		}
	}()
	wg.Wait()
}


func (a *App) getTaskProgress(c *gin.Context) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()
	key := rediskey.CORE_SUBMISSION_PROGRESS.Format(c.Param("sid"))
	lastId := c.Request.Header.Get("Last-Event-Id")
	if lastId == "" {
		lastId = "0"
	}
	c.Header("X-Accel-Buffering", "no") // need this to make SSE work through nginx
	c.Render(200, sse.Event{Event: "start"})
	ping := time.NewTicker(time.Second * 30)
	defer ping.Stop()
	chMsg, chErr := a.Redis.ReadStream(ctx, key, lastId)
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-chErr:
			c.AbortWithError(500, err)
			return
		case msg := <-chMsg:
			var ev sse.Event
			ev.Id = msg.ID
			ev.Event = msg.Data["type"]
			if d, ok := msg.Data["data"]; ok {
				ev.Data = json.RawMessage(d)
			}
			c.Render(-1, ev)
			c.Writer.Flush()
		case <-ping.C:
			c.Render(-1, sse.Event{Event: "ping"})
			c.Writer.Flush()
		}
	}
}
