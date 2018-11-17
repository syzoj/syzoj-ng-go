package main

import (
    "log"
    "flag"
    "encoding/json"
    "io/ioutil"
    _ "github.com/lib/pq"
    "github.com/go-redis/redis"

    "github.com/syzoj/syzoj-ng-go/app"
)

type syzoj_config struct {
    Database string `json:"database"`
    Redis redis.Options `json:"redis"`
    Addr string `json:"addr"`
    GitPath string `json:"git_path"`
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

    err = app_instance.SetupRedis(&config.Redis)
    if err != nil {
        log.Fatal("Error setting up redis:", err)
    }

    err = app_instance.SetupHttpServer(config.Addr)
    if err != nil {
        log.Fatal("Error setting up http server:", err)
    }

    err = app_instance.SetupGitServer(config.GitPath)
    if err != nil {
        log.Fatal("Error setting up git server:", err)
    }
    app_instance.AddGitServer()

    err = app_instance.SetupApiServer()
    if err != nil {
        log.Fatal("Error setting up api server:", err)
    }
    app_instance.AddApiServer()

    app_instance.Run()
}
