package app

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/microcosm-cc/bluemonday"
	"github.com/syzoj/syzoj-ng-go/svc/app/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/russross/blackfriday.v2"
)

type GetProblemsResponse struct {
	Problems []*GetProblemsResponseProblem
	Sqlcnt   int
}

type GetProblemsResponseProblem struct {
	Title null.String `json:"title"`
	Tags  []string    `json:"tags"`
}

func (a *App) getApiProblems(c *gin.Context) {
	ctx := c.Request.Context()
	sqlcnt := 0
	resp := &GetProblemsResponse{Problems: []*GetProblemsResponseProblem{}}
	/*
		problems, err := models.Problems(qm.Where("is_public=1"), qm.OrderBy("id ASC"), qm.Limit(50)).All(ctx, a.Db)
		sqlcnt++
		if err != nil {
			c.AbortWithError(500, err)
			return
		}
		for _, problem := range problems {
			prob := &GetProblemsResponseProblem{Tags: []string{}}
			prob.Title = problem.Title
			tagids, err := models.ProblemTagMaps(qm.Select("tag_id"), qm.Where("problem_id=?", problem.ID)).All(ctx, a.Db)
			sqlcnt++
			if err != nil {
				log.WithError(err).Error("failed to execute SQL")
				continue
			}
			for _, tagid := range tagids {
				tag, err := models.ProblemTags(qm.Select("name", "color"), qm.Where("id=?", tagid.TagID)).One(ctx, a.Db)
				if err != nil {
					log.WithError(err).Error("failed to execute SQL")
					continue
				}
				if tag.Name.Valid {
					prob.Tags = append(prob.Tags, tag.Name.String)
				}
				sqlcnt++
			}
			resp.Problems = append(resp.Problems, prob)
		}
	*/
	aa, err := models.Problems(qm.OrderBy("id ASC")).All(ctx, a.Db)
	bb, err := models.ProblemTagMaps(qm.Select("tag_id")).All(ctx, a.Db)
	cc, err := models.ProblemTags(qm.Select("name", "color")).All(ctx, a.Db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	sqlcnt = len(aa) + len(bb) + len(cc)
	resp.Sqlcnt = sqlcnt
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

	body := &strings.Builder{}
	resp := &GetProblemResponse{}
	resp.Title = problem.Title.String
	if problem.Description.String != "" {
		parseMarkdown(body, "# 题目描述\n\n"+problem.Description.String)
	}
	if problem.InputFormat.String != "" {
		parseMarkdown(body, "# 输入格式\n\n"+problem.InputFormat.String)
	}
	if problem.OutputFormat.String != "" {
		parseMarkdown(body, "# 输出格式\n\n"+problem.OutputFormat.String)
	}
	if problem.Example.String != "" {
		parseMarkdown(body, "# 样例\n\n"+problem.Example.String)
	}
	if problem.LimitAndHint.String != "" {
		parseMarkdown(body, "# 数据范围与提示\n\n"+problem.LimitAndHint.String)
	}
	resp.Body = body.String()
	c.JSON(200, resp)
}

func parseMarkdown(builder *strings.Builder, body string) {
	builder.Write(bluemonday.UGCPolicy().SanitizeBytes(blackfriday.Run([]byte(body))))
}
