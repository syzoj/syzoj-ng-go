package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/config"
	"github.com/syzoj/syzoj-ng-go/lib/life"
	"github.com/syzoj/syzoj-ng-go/svc/app"
	srvredis "github.com/syzoj/syzoj-ng-go/svc/redis"
)

var log = logrus.StandardLogger()

func main() {
	listenAddr := config.GetHttpListenAddr()
	redis, err := config.NewRedis("")
	if err != nil {
		log.WithError(err).Error("failed to get redis config")
		return
	}
	judgeToken := os.Getenv("JUDGE_TOKEN")
	r := srvredis.DefaultRedisService(redis)
	db, err := config.NewMySQL("")
	if err != nil {
		log.WithError(err).Error("failed to get mysql config")
		return
	}
	minio, err := config.NewMinio("")
	if err != nil {
		log.WithError(err).Error("failed to get s3 config")
		return
	}
	_, _ = r, minio
	a := &app.App{
		Db:         db,
		ListenAddr: listenAddr,
		JudgeToken: judgeToken,
		Redis:          r,
		//Minio:          minio,
		//TestdataBucket: "testdata",
	}
	a.Run(life.SignalContext())
}
