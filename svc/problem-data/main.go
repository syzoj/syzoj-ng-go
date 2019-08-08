package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/automation"
	"github.com/syzoj/syzoj-ng-go/lib/config"
	"github.com/valyala/fasthttp"
)

var log = logrus.StandardLogger()

type App struct {
	httpCli       *fasthttp.Client
	automationCli *automation.Client
	dataPath      string
	listenPort    int
	locks         map[string]string
	locksMu       sync.Mutex
}

func (app *App) run() {
	app.locks = make(map[string]string)
	app.dataPath = os.Getenv("DATA_PATH")
	if app.dataPath == "" {
		log.Error("DATA_PATH is empty")
		return
	}
	if !filepath.IsAbs(app.dataPath) {
		log.Warning("DATA_PATH is not absolute")
	}

	automationUrl, err := config.GetHttpURL("AUTOMATION")
	if err != nil {
		log.WithError(err).Error("Failed to get AUTOMATION url")
		return
	}
	app.httpCli = &fasthttp.Client{}
	app.automationCli = automation.NewClient(automationUrl, app.httpCli)

	listenPort, err := config.GetHttpListenPort()
	if err != nil {
		log.WithError(err).Error("Failed to get http listen port")
		return
	}
	app.listenPort = listenPort

	router := fasthttprouter.New()
	router.DELETE("/problem/:name", app.deleteProblem)
	router.PUT("/problem/:name/data", app.putProblemData)
	router.POST("/problem/:name/extract", app.postProblemExtract)
	router.GET("/problem/:name/parse-data", app.getProblemParseData)
	server := &fasthttp.Server{}
	server.Handler = router.Handler
	server.ListenAndServe(fmt.Sprintf(":%d", app.listenPort))
}

func main() {
	app := &App{}
	app.run()
}
