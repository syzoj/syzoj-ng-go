package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgraph-io/dgo"
	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/syzoj/syzoj-ng-go/app/api"
	"github.com/syzoj/syzoj-ng-go/app/judge"
	judge_api "github.com/syzoj/syzoj-ng-go/app/judge/protos"
)

var log = logrus.StandardLogger()

type syzoj_config struct {
	Dgraph string `json:"dgraph"`
	Addr   string `json:"addr"`
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	// Parse config
	configPtr := flag.String("config", "config.json", "Sets the config file")

	flag.Parse()

	var err error
	var configData []byte
	if configData, err = ioutil.ReadFile(*configPtr); err != nil {
		log.Fatal("Error reading config file: ", err)
	}
	var config *syzoj_config
	if err = json.Unmarshal(configData, &config); err != nil {
		log.Fatal("Error parsing config file: ", err)
	}

	log.Info("Connecting to Dgraph")
	var dgraphConn *grpc.ClientConn
	if dgraphConn, err = grpc.Dial(config.Dgraph, grpc.WithInsecure()); err != nil {
		log.Fatal("Error connecting to Dgraph: ", err)
	}
	var dgraph = dgo.NewDgraphClient(dgo_api.NewDgraphClient(dgraphConn))
	defer func() {
		log.Info("Disconnecting from Dgraph")
		dgraphConn.Close()
	}()

	var grpcServer *grpc.Server = grpc.NewServer()

	log.Info("Start judge service")
	var judgeService judge.Service
	if judgeService, err = judge.CreateJudgeService(dgraph); err != nil {
		log.Fatal("Error starting judge service: ", err)
	}
	defer func() {
		log.Info("Stopping judge service")
		judgeService.Close()
	}()
	judge_api.RegisterJudgeServer(grpcServer, judgeService)

	reflection.Register(grpcServer)
	go func() {
		log.Info("Setting up gRPC service")
		lis, err := net.Listen("tcp", "0.0.0.0:3073")
		if err != nil {
			log.Fatal("Failed to listen: ", err)
		}
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC service: ", err)
		}
	}()

	log.Info("Setting up api server")
	var apiServer *api.ApiServer
	if apiServer, err = api.CreateApiServer(dgraph, judgeService); err != nil {
		log.Fatal("Error intializing api server: ", err)
	}

	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(apiServer)
	//router.Handle("/judge-traditional", tjudgeService)
	router.Handle("/", http.FileServer(http.Dir("static")))

	server := &http.Server{
		Addr:         config.Addr,
		Handler:      router,
		WriteTimeout: time.Second * 10,
	}
	go func() {
		log.Infof("Starting web server at %s", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("Web server failed unexpectedly: ", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan
	server.Shutdown(context.Background())
}
