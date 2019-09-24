package app

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/syzoj/syzoj-ng-go/lib/divine"
	"github.com/syzoj/syzoj-ng-go/models"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

type IndexResponse struct {
	Divine   *divine.Divine
	Ranklist []*IndexResponseUser    `json:"ranklist"`
	Notices  []*IndexResponseNotice  `json:"notices"`
	Contests []*IndexResponseContest `json:"contests"`
	Problems []*IndexResponseProblem `json:"problems"`
}
type IndexResponseUser struct {
	ID          int         `json:"id"`
	Username    null.String `json:"username"`
	Nameplate   null.String `json:"nameplate"`
	Information null.String `json:"information"`
}
type IndexResponseContest struct {
	ID        int      `json:"id"`
	StartTime null.Int `json:"start_time"`
	EndTime   null.Int `json:"end_time"`
}
type IndexResponseNotice struct {
	ID    int         `json:"id"`
	Title null.String `json:"title"`
	Date  null.Int    `json:"date"`
}
type IndexResponseProblem struct {
	ID            int         `json:"id"`
	Title         null.String `json:"title"`
	PublicizeTime null.Time   `json:"publicize_time"`
}

func (a *App) getApiIndex(c *gin.Context) {
	ctx := c.Request.Context()
	users, err := models.Users(qm.OrderBy("rating desc"), qm.Limit(20)).All(ctx, a.Db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	articles, err := models.Articles(qm.Where("is_notice = true"), qm.OrderBy("public_time desc")).All(ctx, a.Db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	contests, err := models.Contests(qm.Where("is_public = true"), qm.OrderBy("start_time desc"), qm.Limit(5)).All(ctx, a.Db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	problems, err := models.Problems(qm.Where("is_public = true"), qm.OrderBy("publicize_time desc"), qm.Limit(5)).All(ctx, a.Db)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	resp := &IndexResponse{
		Ranklist: []*IndexResponseUser{},
		Notices:  []*IndexResponseNotice{},
		Contests: []*IndexResponseContest{},
		Problems: []*IndexResponseProblem{},
	}
	for _, user := range users {
		resp.Ranklist = append(resp.Ranklist, &IndexResponseUser{
			ID:          user.ID,
			Username:    user.Username,
			Nameplate:   user.Nameplate,
			Information: user.Information,
		})
	}
	for _, article := range articles {
		resp.Notices = append(resp.Notices, &IndexResponseNotice{
			ID:    article.ID,
			Title: article.Title,
			Date:  article.PublicTime,
		})
	}
	for _, contest := range contests {
		resp.Contests = append(resp.Contests, &IndexResponseContest{
			ID:        contest.ID,
			StartTime: contest.StartTime,
			EndTime:   contest.EndTime,
		})
		_ = contest
	}
	for _, problem := range problems {
		resp.Problems = append(resp.Problems, &IndexResponseProblem{
			ID:            problem.ID,
			Title:         problem.Title,
			PublicizeTime: problem.PublicizeTime,
		})
	}
	/* TODO: fortune */
	userId := c.GetInt(GIN_USER_ID)
	if userId != 0 {
		// TODO: supply sex
		resp.Divine = divine.DoDivine(strconv.Itoa(userId), 0)
	}
	c.JSON(200, resp)
}
