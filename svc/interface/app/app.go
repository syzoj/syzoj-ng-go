// The interface service.
package app

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go"
	"github.com/sirupsen/logrus"
	lredis "github.com/syzoj/syzoj-ng-go/lib/redis"
	"github.com/syzoj/syzoj-ng-go/svc/interface/auth"
	"github.com/syzoj/syzoj-ng-go/svc/interface/problem"
	"github.com/syzoj/syzoj-ng-go/svc/interface/redisscan"
	"github.com/syzoj/syzoj-ng-go/svc/interface/stats"
)

var log = logrus.StandardLogger()

// App represents an 'interface' service.
type App struct {
	ListenAddr     string
	RedisSess      *lredis.PoolWrapper
	RedisStats     *lredis.PoolWrapper
	RedisCache     *lredis.PoolWrapper
	Db             *sqlx.DB
	Minio          *minio.Client
	TestdataBucket string

	ctx       context.Context
	auth      *auth.AuthMiddleware
	prob      *problem.ProblemService
	stats     *stats.Stats
	redisscan *redisscan.Redisscan
}

// Run app
func (app *App) Run() {
	app.ctx = context.Background()
	app.auth = auth.DefaultAuthMiddleware(app.RedisSess)
	app.prob = problem.DefaultProblemService(app.Db, app.RedisSess, app.Minio, app.TestdataBucket)
	app.redisscan = redisscan.DefaultRedisscan(app.RedisSess, app)
	app.stats = &stats.Stats{
		Redis:           app.RedisStats,
		KeyPrefix:       "stats:",
		UpstreamCounter: app.saveCounter,
	}
	router := gin.Default()
	router.Use(app.auth.Handle)
	loginRequired := router.Group("/")
	loginRequired.Use(func(c *gin.Context) {
		inf := app.auth.GetInfo(c)
		if inf.UserId == "" {
			c.JSON(403, gin.H{
				"success": false,
				"message": "Not logged in",
			})
			c.Abort()
		}
	})
	router.POST("/api/register", app.HandleRegister)
	router.POST("/api/login", app.HandleLogin)
	loginRequired.POST("/api/problem/new", app.HandleProblemNew)
	router.GET("/api/problem/id/:problem_id/statement", app.GetProblemStatement)
	router.GET("/api/problem/id/:problem_id", app.GetProblemStatement)
	router.GET("/api/problem/id/:problem_id/judge-info", app.GetProblemJudgeInfo)
	loginRequired.PUT("/api/problem/id/:problem_id", app.HandleProblemPut)
	loginRequired.POST("/api/problem/upload-temp", app.HandleProblemUploadTemp)
	loginRequired.POST("/api/problem/id/:problem_id/upload-data", app.HandleProblemUploadData)
	loginRequired.POST("/api/problem/id/:problem_id/delete-data", app.HandleProblemDeleteData)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		router.Run(app.ListenAddr)
	}()
	go func() {
		defer wg.Done()
		app.redisscan.Run(context.Background())
	}()
	wg.Wait()
}
