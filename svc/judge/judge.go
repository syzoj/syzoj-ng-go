// Judge provides the judging infrastructure.
package judge

import (
	"context"
	"database/sql"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/rediskey"
	svcredis "github.com/syzoj/syzoj-ng-go/svc/redis"
	"github.com/volatiletech/null"
)

var log = logrus.StandardLogger()

type JudgeService struct {
	Db    *sql.DB
	Redis *svcredis.RedisService
}

func DefaultJudgeService(db *sql.DB, r *svcredis.RedisService) *JudgeService {
	return &JudgeService{Db: db, Redis: r}
}

func (s *JudgeService) Run(ctx context.Context) error {
	_, err := s.Redis.DoContext(ctx, "XGROUP", "CREATE", rediskey.CORE_QUEUE.Format("default"), "judger", "$", "MKSTREAM")
	if err != nil {
		return err
	}
	return nil
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
type SubmissionResultSummary struct {
	Status string
	Time   null.Float64
	Memory null.Float64
	Score  null.Float64
}

func GetSummary(res *Judge) *SubmissionResultSummary {
	var (
		status string
		rtime  null.Float64
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
		rtime = null.Float64From(0)
		memory = null.Float64From(0)
		for _, subtask := range prog.Judge.Subtasks {
			if subtask == nil {
				continue
			}
			for _, c := range subtask.Cases {
				if c.Result != nil {
					rtime.Float64 += c.Result.Time.Float64
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
	}
	return &SubmissionResultSummary{
		Status: status,
		Time:   rtime,
		Memory: memory,
		Score:  score,
	}
}

func (s *JudgeService) SaveTask(ctx context.Context, sid string, res *Judge) error {

	// Notify callbacks
	keys, err := redis.Values(s.Redis.DoContext(ctx, "SMEMBERS", rediskey.CORE_SUBMISSION_CALLBACK.Format(sid)))
	if err != nil {
		return err
	}
	pipe, err := s.Redis.NewPipeline(ctx)
	if err != nil {
		return err
	}
	defer pipe.Close()
	for _, key := range keys {
		s, err := redis.String(key, nil)
		if err != nil {
			return err
		}
		pipe.Do(nil, "XADD", s, "*", "sid", sid)
	}
	// Make submission expire
	pipe.Do(nil, "EXPIRE", rediskey.CORE_SUBMISSION_PROGRESS.Format(sid), int64(rediskey.DEFAULT_EXPIRE/time.Second))
	pipe.Do(nil, "EXPIRE", rediskey.CORE_SUBMISSION_DATA.Format(sid), int64(rediskey.DEFAULT_EXPIRE/time.Second))
	pipe.Do(nil, "EXPIRE", rediskey.CORE_SUBMISSION_CALLBACK.Format(sid), int64(rediskey.DEFAULT_EXPIRE/time.Second))
	return pipe.Flush(ctx)
}
