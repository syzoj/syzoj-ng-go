package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/syzoj/syzoj-ng-go/database"
	"github.com/syzoj/syzoj-ng-go/server"
	"github.com/syzoj/syzoj-ng-go/server/handlers"
)

var log = logrus.StandardLogger()

type syzoj_config struct {
	MySQL   string              `json:"mysql"`
	Addr    string              `json:"addr"`
	RpcAddr string              `json:"rpc_addr"`
	Server  server.ServerConfig `json:"server"`
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Must specify a subcommand")
		return
	}
	switch os.Args[1] {
	case "run":
		cmdRun()
	default:
		fmt.Println("Must specify a subcommand")
		return
	}
}

func cmdRun() {
	runFlagSet := flag.NewFlagSet("run", flag.ExitOnError)
	configPtr := runFlagSet.String("config", "config.json", "Sets the config file")
	runFlagSet.Parse(os.Args[2:])
	log.SetLevel(logrus.DebugLevel)

	var err error
	var configData []byte
	if configData, err = ioutil.ReadFile(*configPtr); err != nil {
		log.Fatal("Error reading config file: ", err)
	}
	var config *syzoj_config
	if err = json.Unmarshal(configData, &config); err != nil {
		log.Fatal("Error parsing config file: ", err)
	}

	log.Info("Connecting to MySQL")
	var db *database.Database
	if db, err = database.Open("mysql", config.MySQL); err != nil {
		log.Fatal("Error connecting to MySQL: ", err)
	}
	defer func() {
		log.Info("Disconnecting from MySQL")
		db.Close()
	}()

	var grpcServer *grpc.Server = grpc.NewServer()

	log.Info("Start SYZOJ")
	var s *server.Server
	s = server.NewServer(db, &config.Server)
	handlers.RegisterHandlers(s.ApiServer())
	defer func() {
		log.Info("Stopping SYZOJ")
		s.Close()
	}()
	reflection.Register(grpcServer)
	go func() {
		log.Info("Setting up gRPC service")
		lis, err := net.Listen("tcp", config.RpcAddr)
		if err != nil {
			log.Fatal("Failed to listen: ", err)
		}
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatal("Failed to serve gRPC service: ", err)
		}
	}()

	log.Info("Setting up HTTP server")
	router := mux.NewRouter()
	router.PathPrefix("/api").Handler(s.ApiServer())
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
	log.Info("Shutting down web server")
	server.Close()
}
