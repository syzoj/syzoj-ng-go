package app

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/models"
	"github.com/syzoj/syzoj-ng-go/svc/judge"
	svcredis "github.com/syzoj/syzoj-ng-go/svc/redis"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

var log = logrus.StandardLogger()

const GIN_USER_ID = "USER_ID"

type App struct {
	Db           *sql.DB
	ListenAddr   string
	Redis        *svcredis.RedisService // The persistent redis instance. No eviction policies allowed.
	RedisCache   *svcredis.RedisService
	JudgeService *judge.JudgeService
	JudgeToken   string
}

func DefaultApp(db *sql.DB, redis *svcredis.RedisService, redisCache *svcredis.RedisService, listenAddr string, judgeService *judge.JudgeService) *App {
	return &App{
		Db:           db,
		ListenAddr:   listenAddr,
		Redis:        redis,
		RedisCache:   redisCache,
		JudgeService: judgeService,
	}
}

func (a *App) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	router.GET("/api/submission-progress/:sid", a.getTaskProgress)
	router.GET("/api/header", a.getHeader)
	jg := router.Group("/judge")
	jg.Use(a.useCheckJudgeToken)
	jg.GET("/wait-for-task", a.getJudgeWaitForTask)
	server := &http.Server{Addr: a.ListenAddr, Handler: router}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		server.Close()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil {
			log.WithError(err).Error("failed to listen and serve")
		}
	}()
	go func() {
		defer wg.Done()
		if err := a.handleJudgeDone(ctx); err != nil {
			log.WithError(err).Error("failed to handle judge done")
		}
	}()
	wg.Wait()
	return nil
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
