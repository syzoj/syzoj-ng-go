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
	"github.com/syzoj/syzoj-ng-go/svc/app/models"
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
				_, err := a.Redis.DoContext(ctx, "XADD", rediskey.CORE_SUBMISSION_PROGRESS.Format(sid), "*", "type", "reset")
				if err != nil {
					log.WithError(err).Error("failed to call XADD")
					return
				}
				if err := conn.WriteJSON(gin.H{
					"id":   msg.ID,
					"sid":  sid,
					"data": json.RawMessage(msg.Data["data"]),
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
				res := &Judge{}
				if err := json.Unmarshal(data.Data, res); err != nil {
					log.WithField("sid", data.Sid).WithError(err).Error("failed to parse data")
					return
				}
				if res.Type == null.IntFrom(4) {
					if err := a.saveSubmission(ctx, data.Sid, res); err != nil {
						log.WithField("sid", data.Sid).WithError(err).Error("failed to save submission")
						return
					}
					if _, err := a.Redis.DoContext(ctx, "XADD", skey, "*", "type", "done"); err != nil {
						log.WithError(err).Error("failed to call XADD")
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

// ConvertResult
type Judge struct {
	TaskId   null.String    `json:"task_id,omitempty"`
	Type     null.Int       `json:"type,omitempty"`
	Progress *JudgeProgress `json:"progress,omitempty"`
}
type JudgeProgress struct {
	Compile       *JudgeProgressCompile `json:"compile,omitempty"`
	Judge         *JudgeProgressJudge   `json:"judge,omitempty"`
	Error         null.Int              `json:"error,omitempty"`
	SystemMessage null.String           `json:"systemMessage,omitempty"`
	Status        null.Int              `json;"status,omitempty"`
	Message       null.String           `json:"message,omitempty"`
}
type JudgeProgressCompile struct {
	Status null.Int `json:"status"`
}
type JudgeProgressJudge struct {
	Subtasks []*JudgeProgressJudgeSubtask `json:"subtasks,omitempty"`
}
type JudgeProgressJudgeSubtask struct {
	Type  null.Int                         `json:"type,omitempty"`
	Cases []*JudgeProgressJudgeSubtaskCase `json:"cases,omitempty"`
	Score null.Float64                     `json:"score,omitempty"`
}
type JudgeProgressJudgeSubtaskCase struct {
	Status       null.Int                 `json:"status,omitempty"`
	Result       *JudgeProgressCaseResult `json:"result,omitempty"`
	ErrorMessage null.String              `json:"errorMessage,omitempty"`
}
type JudgeProgressCaseResult struct {
	Type          null.Int     `json:"type,omitempty"`
	Time          null.Float64 `json:"time,omitempty"`
	Memory        null.Float64 `json:"memory,omitempty"`
	Status        null.Int     `json:"status,omitempty"`
	Input         *FileContent `json:"input,omitempty"`
	Output        *FileContent `json:"output,omitempty"`
	Score         null.Float64 `json:"score,omitempty"`
	ScoringRate   null.Float64 `json:"scoringRate,omitempty"` // ???
	UserOutput    null.String  `json:"userOutput,omitempty'`
	UserError     null.String  `json:"userError,omitempty"`
	SpjMessage    null.String  `json:"spjMessage,omitempty"`
	SystemMessage null.String  `json:"systemMessage,omitempty"`
}
type FileContent struct {
	Content string `json:"content"`
	Name    string `json:"name"`
}
type SubmissionResult struct {
	Score        null.Int     `json:"score,omitempty"`
	Pending      null.Bool    `json:"pending,omitempty"`
	StatusString null.String  `json:"statusString,omitempty"`
	Time         null.Float64 `json:"time,omitempty"`
	Memory       null.Float64 `json:"memory,omitempty"`
	Result       null.String  `json:"result,omitempty"`
}

var statusString = map[int]string{
	1:  "Accepted",
	2:  "Wrong Answer",
	3:  "Partially Correct",
	4:  "Memory Limit Exceeded",
	5:  "Time Limit Exceeded",
	6:  "Output Limit Exceeded",
	7:  "Runtime Error",
	8:  "File Error",
	9:  "Judgement Failed",
	10: "Invalid Interaction",
}

// TODO: updateRelatedInfo
func (a *App) saveSubmission(ctx context.Context, sid string, res *Judge) error {
	var (
		status string
		time   null.Float64
		memory null.Float64
		score  null.Float64
	)
	prog := res.Progress
	if prog.Compile != nil && prog.Compile.Status == null.IntFrom(3) {
		status = "Compile Error"
	} else if prog.Error.Valid {
		switch prog.Error.Int {
		case 0:
			status = "No Testdata"
		default:
			status = "System Error"
		}
	} else if prog.Judge != nil && prog.Judge.Subtasks != nil {
		time = null.Float64From(0)
		memory = null.Float64From(0)
		for _, subtask := range prog.Judge.Subtasks {
			if subtask == nil {
				continue
			}
			for _, c := range subtask.Cases {
				if c.Result != nil {
					time.Float64 += c.Result.Time.Float64
					if c.Result.Memory.Valid && memory.Float64 < c.Result.Memory.Float64 {
						memory.Float64 = c.Result.Memory.Float64
					}
					if status == "" && c.Result.Type != null.IntFrom(1) {
						status = statusString[c.Result.Type.Int]
					}
				}
			}
			score.Float64 += subtask.Score.Float64
		}
		if status == "" {
			status = statusString[1]
		}
	} else {
		status = "System Error"
		log.Infof("no subtasks system error: %#v", prog)
	}
	subm, err := models.JudgeStates(qm.Where("task_id=?", sid)).One(ctx, a.Db)
	if err != nil {
		return err
	}
	subm.Score = null.IntFrom(int(score.Float64))
	subm.Pending = null.Int8From(0)
	subm.Status = null.StringFrom(status)
	subm.TotalTime = null.IntFrom(int(time.Float64))
	subm.MaxMemory = null.IntFrom(int(memory.Float64))
	progBytes, err := json.Marshal(prog)
	if err != nil {
		return err
	}
	subm.Result = null.StringFrom(string(progBytes))
	if _, err := subm.Update(ctx, a.Db, boil.Blacklist()); err != nil {
		return err
	}
	return nil
}

func (a *App) getSubmissionProgress(c *gin.Context) {
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
