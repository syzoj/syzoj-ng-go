package app

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/syzoj/syzoj-ng-go/lib/rediskey"
	"github.com/syzoj/syzoj-ng-go/models"
	"github.com/syzoj/syzoj-ng-go/svc/judge"
	svcredis "github.com/syzoj/syzoj-ng-go/svc/redis"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
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
	pipeline, err := a.Redis.NewPipeline(ctx)
	if err != nil {
		log.WithError(err).Error("failed to create redis pipeline")
		return
	}
	defer pipeline.Close()

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
				pipeline.Do(a.checkRedis, "XADD", rediskey.CORE_SUBMISSION_PROGRESS.Format(sid), "*", "type", "reset")
				if err := pipeline.Flush(ctx); err != nil {
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
				pipeline.Do(a.checkRedis, "XADD", skey, "*", "type", "progress", "data", []byte(data.Data))
			case "finish":
				res := &judge.Judge{}
				if err := json.Unmarshal(data.Data, res); err != nil {
					log.WithField("sid", data.Sid).WithError(err).Error("failed to parse data")
					return
				}
				if res.Type == null.IntFrom(4) {
					// Save submission result
					sum := judge.GetSummary(res)
					subm, err := models.JudgeStates(qm.Where("task_id=?", data.Sid)).One(ctx, a.Db)
					if err != nil {
						log.WithError(err).Error("failed to query db")
						return
					}
					subm.Score = null.IntFrom(int(sum.Score.Float64))
					subm.Pending = null.Int8From(0)
					subm.Status = null.StringFrom(sum.Status)
					subm.TotalTime = null.IntFrom(int(sum.Time.Float64))
					subm.MaxMemory = null.IntFrom(int(sum.Memory.Float64))
					progBytes, err := json.Marshal(res.Progress)
					if err != nil {
						panic(err)
					}
					subm.Result = null.StringFrom(string(progBytes))
					if _, err := subm.Update(ctx, a.Db, boil.Blacklist()); err != nil {
						log.WithError(err).Error("failed to update db")
						return
					}
					pipeline.Do(a.checkRedis, "SET", rediskey.CORE_SUBMISSION_RESULT.Format(data.Sid), data.Data)
					if err := pipeline.Flush(ctx); err != nil {
						log.WithError(err).Error("failed to send redis")
						return
					}
					if err := a.JudgeService.SaveTask(ctx, data.Sid, res); err != nil {
						log.WithField("sid", data.Sid).WithError(err).Error("failed to save submission")
						return
					}
					pipeline.Do(a.checkRedis, "XADD", skey, "*", "type", "done")
					var res svcredis.RedisResult
					pipeline.Do(res.Save, "XACK", key, "judger", data.Id)
					if err := pipeline.Flush(ctx); err != nil {
						log.WithError(err).Error("failed to flush to Redis")
						return
					}
					n, err := redis.Int64(res.Result, res.Err)
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

// Handle judge done events and update stats.
func (a *App) handleJudgeDone(ctx context.Context) error {
	pipe, err := a.Redis.NewPipeline(ctx)
	if err != nil {
		return err
	}
	pipe.Do(nil, "XGROUP", "CREATE", rediskey.MAIN_JUDGE_DONE, "main", "$", "MKSTREAM")
	if err := pipe.Flush(ctx); err != nil {
		return err
	}
	sema := make(chan struct{}, 1)
	sema <- struct{}{}
	chMsg, chErr := a.Redis.ReadStreamGroup(ctx, rediskey.MAIN_JUDGE_DONE, "main", "main", sema)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				pipe.Do(nil, "PING")
				if err := pipe.Flush(ctx); err != nil {
					log.WithError(err).Error("failed to flush to redis")
					return
				}
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-chErr:
			return err
		case msg := <-chMsg:
			sid := msg.Data["sid"]
			if sid == "" {
				log.WithField("id", msg.ID).Error("judge done stream: missing sid field")
				continue
			}
			subm, err := models.JudgeStates(qm.Where("task_id=?", sid)).One(ctx, a.Db)
			if err != nil {
				log.WithField("sid", sid).Error("judge done stream: failed to get submission")
				continue
			}
			if subm.ProblemID.Valid {
				pid := strconv.Itoa(subm.ProblemID.Int)
				pipe.Do(a.checkRedis, "INCR", rediskey.MAIN_PROBLEM_SUBMITS.Format(pid))
				if subm.Status == null.StringFrom("Accepted") {
					pipe.Do(a.checkRedis, "INCR", rediskey.MAIN_PROBLEM_ACCEPTS.Format(pid))
				}
				if subm.UserID.Valid {
					uid := strconv.Itoa(subm.UserID.Int)
					pipe.Do(a.checkRedis, "HSET", rediskey.MAIN_USER_LAST_SUBMISSION.Format(uid), pid, sid)
					if subm.Status == null.StringFrom("Accepted") {
						pipe.Do(a.checkRedis, "HSET", rediskey.MAIN_USER_LAST_ACCEPT.Format(uid), pid, sid)
					}
				}
			}
			pipe.Do(a.checkRedis, "XACK", rediskey.MAIN_JUDGE_DONE, "main", msg.ID)
		}
	}
}
