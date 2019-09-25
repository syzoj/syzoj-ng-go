package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/microcosm-cc/bluemonday"
	"github.com/syzoj/syzoj-ng-go/lib/rediskey"
	"github.com/syzoj/syzoj-ng-go/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/russross/blackfriday.v2"
)

type GetProblemsResponse struct {
	Problems []*GetProblemsResponseProblem
}

type GetProblemsResponseProblem struct {
	Id         string      `json:"id"`
	Title      null.String `json:"title"`
	Tags       []string    `json:"tags"`
	SubmitNum  int64       `json:"submit_num"`
	AcceptNum  int64       `json:"accept_num"`
	LastSubmit null.String `json:"last_submit"`
	LastAccept null.String `json:"last_accept"`
}

func (a *App) getApiProblems(c *gin.Context) {
	ctx := c.Request.Context()
	resp := &GetProblemsResponse{Problems: []*GetProblemsResponseProblem{}}
	problems, err := models.Problems(qm.Where("is_public=1"), qm.OrderBy("id ASC")).All(ctx, a.Db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	pipeline, err := a.Redis.NewPipeline(ctx)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	defer pipeline.Close()
	userIdInt := c.GetInt(GIN_USER_ID)
	var userId string
	if userIdInt != 0 {
		userId = strconv.Itoa(userIdInt)
	}
	for _, problem := range problems {
		prob := &GetProblemsResponseProblem{Tags: []string{}}
		prob.Title = problem.Title
		if problem.Tags.Valid {
			json.Unmarshal([]byte(problem.Tags.String), &prob.Tags)
		}
		resp.Problems = append(resp.Problems, prob)
		probId := strconv.Itoa(problem.ID)
		prob.Id = probId
		pipeline.Do(func(data interface{}, err error) {
			num, err := redis.Int64(data, err)
			if err != nil {
				log.WithField("problem_id", probId).WithError(err).Error("failed to get problem submit count")
				return
			}
			prob.SubmitNum = num
		}, "GET", rediskey.MAIN_PROBLEM_SUBMITS.Format(probId))
		pipeline.Do(func(data interface{}, err error) {
			num, err := redis.Int64(data, err)
			if err != nil {
				log.WithField("problem_id", probId).WithError(err).Error("failed to get problem submit count")
				return
			}
			prob.AcceptNum = num
		}, "GET", rediskey.MAIN_PROBLEM_ACCEPTS.Format(probId))
		if userId != "" {
			pipeline.Do(func(data interface{}, err error) {
				sid, err := redis.String(data, err)
				if err == redis.ErrNil {
					return
				}
				if err != nil {
					log.WithField("user_id", userId).WithError(err).Error("failed to get user last submission")
					return
				}
				prob.LastSubmit = null.StringFrom(sid)
			}, "HGET", rediskey.MAIN_USER_LAST_SUBMISSION.Format(userId), probId)
			pipeline.Do(func(data interface{}, err error) {
				sid, err := redis.String(data, err)
				if err == redis.ErrNil {
					return
				}
				if err != nil {
					log.WithField("user_id", userId).WithError(err).Error("failed to get user last submission")
					return
				}
				prob.LastAccept = null.StringFrom(sid)
			}, "HGET", rediskey.MAIN_USER_LAST_ACCEPT.Format(userId), probId)
		}
	}
	if err := pipeline.Flush(ctx); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, resp)
}

type GetProblemResponse struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (a *App) getApiProblem(c *gin.Context) {
	ctx := c.Request.Context()
	problemId, err := strconv.Atoi(c.Param("problem_id"))
	if err != nil {
		c.AbortWithError(400, err)
		return
	}
	problem, err := models.Problems(qm.Where("id=?", problemId)).One(ctx, a.Db)
	if err == sql.ErrNoRows {
		c.AbortWithError(404, err)
		return
	}
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if problem.IsPublic.Int8 != 1 {
		c.JSON(200, gin.H{"error": "您没有权限进行此操作。"})
		return
	}

	body := &bytes.Buffer{}
	resp := &GetProblemResponse{}
	resp.Title = problem.Title.String
	if problem.Description.String != "" {
		body.WriteString("# 题目描述\n\n" + problem.Description.String)
	}
	if problem.InputFormat.String != "" {
		body.WriteString("# 输入格式\n\n" + problem.InputFormat.String)
	}
	if problem.OutputFormat.String != "" {
		body.WriteString("# 输出格式\n\n" + problem.OutputFormat.String)
	}
	if problem.Example.String != "" {
		body.WriteString("# 样例\n\n" + problem.Example.String)
	}
	if problem.LimitAndHint.String != "" {
		body.WriteString("# 数据范围与提示\n\n" + problem.LimitAndHint.String)
	}
	resp.Body = string(bluemonday.UGCPolicy().SanitizeBytes(blackfriday.Run(body.Bytes())))
	c.JSON(200, resp)
}
