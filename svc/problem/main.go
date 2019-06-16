package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/elastic/go-elasticsearch"
	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/binlog"
	"github.com/syzoj/syzoj-ng-go/lib/life"
	"github.com/syzoj/syzoj-ng-go/lib/mysql"
	"github.com/valyala/fasthttp"
)

var log = logrus.StandardLogger()

type Main struct {
	dbProb    *sql.DB
	binlogsub *binlog.Subscriber
	es        *elasticsearch.Client
}

func (m *Main) start() error {
	// Connect to MySQL
	dbProb, err := mysql.OpenMySQL("PROBLEM")
	if err != nil {
		return fmt.Errorf("failed to open problem db: %v", err)
	}
	m.dbProb = dbProb
	// Subscribe to binlog
	binlogsub, err := binlog.NewSubscriber("PROBLEM", "problem_events")
	if err != nil {
		log.WithError(err).Error("Failed to create subscriber")
		os.Exit(1)
	}
	m.binlogsub = binlogsub
	// Connect to Elasticsearch
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{os.Getenv("ELASTIC_ADDR")},
	})
	if err != nil {
		log.WithError(err).Error("Failed to connect to Elasticsearch")
		os.Exit(1)
	}
	for {
		_, err := es.Info()
		if err != nil {
			log.WithError(err).Error("Failed to ping Elasticsearch")
			time.Sleep(time.Second)
		} else {
			break
		}
	}
	m.es = es

	ctx := life.SignalContext()
	var wg sync.WaitGroup
	wg.Add(2)
	// Start binlog consumption
	go func() {
		defer wg.Done()
		if err := m.binlogsub.RunPosFile(ctx, m, path.Join(os.Getenv("DATA_PATH"), "POSITION")); err != nil {
			log.WithError(err).Error("Failed to run subscriber")
		}
	}()
	// Start http server
	go func() {
		serv := &fasthttp.Server{}
		router := fasthttprouter.New()
		router.GET("/api/problem/id/:id", m.getProblemId)
		router.POST("/api/problem", m.postProblem)
		router.PUT("/api/problem/id/:id", m.putProblemId)
		router.DELETE("/api/problem/id/:id", m.deleteProblemId)
		router.POST("/api/problem/search", m.postSearch)
		serv.Handler = router.Handler
		go func() {
			<-ctx.Done()
			serv.Shutdown()
			wg.Done()
		}()
		addr := os.Getenv("HTTP_LISTEN_ADDR")
		log.Infof("Starting server at %s", addr)
		if err := serv.ListenAndServe(addr); err != nil {
			log.WithError(err).Error("Failed to run HTTP server")
		}
	}()
	wg.Wait()
	return nil
}

func (m *Main) handleNotFound(ctx *fasthttp.RequestCtx, err error) {
	type resp struct {
		Err string `json:"error"`
	}
	b, _ := json.Marshal(resp{Err: err.Error()})
	ctx.SetStatusCode(404)
	ctx.Success("application/json", b)
}

func (m *Main) handleError(ctx *fasthttp.RequestCtx, err error) {
	log.Errorf("handle %v: err %v", ctx.URI(), err)
	ctx.SetStatusCode(500)
}

func (m *Main) handleBadRequest(ctx *fasthttp.RequestCtx, err error) {
	log.Warningf("handle %v: bad request %v", ctx.URI(), err)
	ctx.SetStatusCode(400)
}

func (m *Main) sendBody(ctx *fasthttp.RequestCtx, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		m.handleError(ctx, err)
		return
	}
	ctx.Success("application/json", b)
}

func main() {
	m := &Main{}
	if err := m.start(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
