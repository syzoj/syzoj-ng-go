package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/config"
	lredis "github.com/syzoj/syzoj-ng-go/lib/redis"
	"github.com/syzoj/syzoj-ng-go/svc/interface/app"
)

var log = logrus.StandardLogger()

func main() {
	listenPort, err := config.GetHttpListenPort()
	if err != nil {
		log.WithError(err).Error("failed to get http listen port")
		return
	}
	redis, err := config.NewRedis("")
	if err != nil {
		log.WithError(err).Error("failed to get redis config")
		return
	}
	redisp := lredis.WrapPool(redis)
	db, err := config.NewMySQLx("")
	if err != nil {
		log.WithError(err).Error("failed to get mysql config")
		return
	}
	minio, err := config.NewMinio("")
	if err != nil {
		log.WithError(err).Error("failed to get s3 config")
		return
	}
	a := &app.App{
		ListenAddr: fmt.Sprintf(":%d", listenPort),
		RedisSess:  redisp,
		RedisStats: redisp,
		Db:         db,
		Minio: minio,
		TestdataBucket: "testdata",
	}
	a.Run()
}
