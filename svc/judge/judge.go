// Judge provides the judging infrastructure.
package judge

import (
	"context"
	"encoding/json"
	"database/sql"

	"github.com/sirupsen/logrus"
	svcredis "github.com/syzoj/syzoj-ng-go/svc/redis"
	"github.com/syzoj/syzoj-ng-go/lib/rediskey"
	"github.com/syzoj/syzoj-ng-go/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/null"
)
var log = logrus.StandardLogger()

type JudgeService struct {
	Db *sql.DB
	Redis *svcredis.RedisService
}

func DefaultJudgeService(db *sql.DB, r *svcredis.RedisService) *JudgeService {
	return &JudgeService{Db: db,Redis: r}
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
func (s *JudgeService) SaveTask(ctx context.Context, sid string, res *Judge) error {
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
	subm, err := models.JudgeStates(qm.Where("task_id=?", sid)).One(ctx, s.Db)
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
	if _, err := subm.Update(ctx, s.Db, boil.Blacklist()); err != nil {
		return err
	}
	return nil
}

