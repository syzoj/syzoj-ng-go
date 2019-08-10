package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/elastic/go-elasticsearch"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/automation"
	"github.com/syzoj/syzoj-ng-go/lib/config"
	"github.com/valyala/fasthttp"
)

var log = logrus.StandardLogger()

type App struct {
	dbUser        *sql.DB
	dbProblem     *sql.DB
	redisSession  *redis.Pool
	redisCache    *redis.Pool
	listenPort    int
	automationCli *automation.Client
	httpCli       *fasthttp.Client
	esProblem     *elasticsearch.Client
}

func (app *App) run() {
	dbUser, err := config.OpenMySQL("USER")
	if err != nil {
		log.WithError(err).Error("Failed to open USER")
		return
	}
	app.dbUser = dbUser

	dbProblem, err := config.OpenMySQL("PROBLEM")
	if err != nil {
		log.WithError(err).Error("Failed to open PROBLEM")
		return
	}
	app.dbProblem = dbProblem

	redisSession, err := config.OpenRedis("SESSION")
	if err != nil {
		log.WithError(err).Error("Failed to open SESSION")
		return
	}
	app.redisSession = redisSession

	redisCache, err := config.OpenRedis("CACHE")
	if err != nil {
		log.WithError(err).Error("Failed to open CACHE")
		return
	}
	app.redisCache = redisCache

	listenPort, err := config.GetHttpListenPort()
	if err != nil {
		log.WithError(err).Error("Failed to get http listen port")
		return
	}
	app.listenPort = listenPort

	automationUrl, err := config.GetHttpURL("AUTOMATION")
	if err != nil {
		log.WithError(err).Error("Failed to get AUTOMATION url")
		return
	}
	app.httpCli = &fasthttp.Client{}
	app.automationCli = automation.NewClient(automationUrl, app.httpCli)

	esProblem, err := config.OpenElastic("PROBLEM")
	if err != nil {
		log.WithError(err).Error("Failed to open elasticsearch for PROBLEM")
		return
	}
	app.esProblem = esProblem

	router := fasthttprouter.New()
	router.POST("/user/register", app.postUserRegister)
	router.POST("/user/login", app.postUserLogin)
	router.GET("/user/current", app.getUserCurrent)
	router.POST("/problem/new", app.postProblemNew)
	router.POST("/problem/id/:uid/upload-data", app.postProblemUploadData)
	router.POST("/problem/id/:uid/submit", app.postProblemSubmit)
	router.GET("/problem/id/:uid/info", app.getProblemInfo)
	router.GET("/problems", app.getProblems)

	server := &fasthttp.Server{
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
	}
	server.Handler = router.Handler
	server.ListenAndServe(fmt.Sprintf(":%d", app.listenPort))
}

func main() {
	app := &App{}
	app.run()
}
