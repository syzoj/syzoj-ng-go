package main

import (
	"database/sql"
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/automation"
	"github.com/syzoj/syzoj-ng-go/lib/config"
	"github.com/valyala/fasthttp"
)

var log = logrus.StandardLogger()

type App struct {
	dbProblem     *sql.DB
	redisSession  *redis.Pool
	httpCli       *fasthttp.Client
	automationCli *automation.Client
	listenPort    int
}

func (app *App) run() {
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

	router := fasthttprouter.New()
	router.POST("/judge-queue/enqueue", app.postEnqueue)
	router.POST("/judge-queue/fetch", app.postFetch)
	router.POST("/judge-queue/task/id/:uid/handle", app.postTaskHandle)
	server := &fasthttp.Server{}
	server.Handler = router.Handler
	server.ListenAndServe(fmt.Sprintf(":%d", app.listenPort))
}

func main() {
	app := &App{}
	app.run()
}
