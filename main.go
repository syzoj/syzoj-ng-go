package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"

	"github.com/syzoj/syzoj-ng-go/app"
)

type syzoj_config struct {
	Database string `json:"database"`
	Addr     string `json:"addr"`
	GitPath  string `json:"git_path"`
	LevelDB  string `json:"leveldb_path"`
}

func main() {
	// Parse config
	configPtr := flag.String("config", "config.json", "Sets the config file")

	flag.Parse()

	var err error
	dat, err := ioutil.ReadFile(*configPtr)
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	var config *syzoj_config
	err = json.Unmarshal(dat, &config)
	if err != nil {
		log.Fatal("Error parsing config file: ", err)
	}

	// Prepare dependencies
	app_instance := &app.App{}
	err = app_instance.SetupDB(config.Database)
	if err != nil {
		log.Fatal("Error setting up database:", err)
	}

	err = app_instance.SetupLevelDB(config.LevelDB)
	if err != nil {
		log.Fatal("Error setting up LevelDB:", err)
	}

	err = app_instance.SetupSessionService()
	if err != nil {
		log.Fatal("Error setting up session service:", err)
	}

	err = app_instance.SetupAuthService()
	if err != nil {
		log.Fatal("Error setting up auth service:", err)
	}

	err = app_instance.SetupProblemsetService()
	if err != nil {
		log.Fatal("Error setting up problemset service:", err)
	}

	// Setup services
	err = app_instance.SetupHttpServer(config.Addr)
	if err != nil {
		log.Fatal("Error setting up http server:", err)
	}

	err = app_instance.SetupGitServer(config.GitPath)
	if err != nil {
		log.Fatal("Error setting up git server:", err)
	}

	err = app_instance.AddGitServer()
	if err != nil {
		log.Fatal("Error adding git server:", err)
	}

	err = app_instance.SetupApiServer()
	if err != nil {
		log.Fatal("Error setting up api server:", err)
	}

	err = app_instance.AddApiServer()
	if err != nil {
		log.Fatal("Error adding api server:", err)
	}

	app_instance.Run()
}
