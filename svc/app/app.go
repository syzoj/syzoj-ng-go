package app

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/svc/app/models"
	svcredis "github.com/syzoj/syzoj-ng-go/svc/redis"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

var log = logrus.StandardLogger()

const GIN_USER_ID = "USER_ID"

type App struct {
	Db         *sql.DB
	ListenAddr string
	Redis      *svcredis.RedisService
	JudgeToken string
}

func DefaultApp(db *sql.DB, redis *svcredis.RedisService, listenAddr string) *App {
	return &App{
		Db:         db,
		ListenAddr: listenAddr,
		Redis:      redis,
	}
}

func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.ensureQueue(ctx, "default"); err != nil {
			log.WithError(err).Error("failed to create default queue")
		}
	}()
	router := gin.Default()
	router.Use(a.UserMiddleware)
	router.GET("/api/index", a.getApiIndex)
	router.POST("/api/login", a.postApiLogin)
	router.GET("/api/problems", a.getApiProblems)
	router.GET("/api/problem/:problem_id", a.getApiProblem)
	router.GET("/api/submission-progress/:sid", a.getSubmissionProgress)
	jg := router.Group("/judge")
	jg.Use(a.useCheckJudgeToken)
	jg.GET("/wait-for-task", a.getJudgeWaitForTask)
	server := &http.Server{Addr: a.ListenAddr, Handler: router}
	go func() {
		<-ctx.Done()
		server.Close()
	}()
	return server.ListenAndServe()
}

func (a *App) UserMiddleware(c *gin.Context) {
	ctx := c.Request.Context()
	c.Set(GIN_USER_ID, 0)
	loginCookie, err := c.Cookie("login")
	if err != nil {
		return
	}
	var loginData []string
	if err := json.Unmarshal([]byte(loginCookie), &loginData); err != nil {
		return
	}
	if len(loginData) != 2 {
		return
	}
	user, err := models.Users(qm.Select("id", "password"), qm.Where("username=?", loginData[0])).One(ctx, a.Db)
	if err != nil {
		return
	}
	if !user.Password.Valid || subtle.ConstantTimeCompare([]byte(loginData[1]), []byte(user.Password.String)) != 1 {
		return
	}
	c.Set(GIN_USER_ID, user.ID)
}
