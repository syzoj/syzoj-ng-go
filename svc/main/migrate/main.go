package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/config"
	"github.com/syzoj/syzoj-ng-go/lib/life"
	"github.com/syzoj/syzoj-ng-go/svc/migrate"
	srvredis "github.com/syzoj/syzoj-ng-go/svc/redis"
)

var log = logrus.StandardLogger()
var funcs = map[string]func(*migrate.MigrateService, context.Context) error{
	"all":              (*migrate.MigrateService).All,
	"problem-tags":     (*migrate.MigrateService).MigrateProblemTags,
	"problem-counter":  (*migrate.MigrateService).MigrateProblemCounter,
	"user-submissions": (*migrate.MigrateService).MigrateUserSubmissions,
}

func main() {
	var flist []func(*migrate.MigrateService, context.Context) error
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if f, ok := funcs[arg]; ok {
			flist = append(flist, f)
		} else {
			log.WithField("arg", arg).Error("unsupported argument")
		}
	}
	redis, err := config.NewRedis("")
	if err != nil {
		log.WithError(err).Error("failed to get redis config")
		return
	}
	r := srvredis.DefaultRedisService(redis)
	db, err := config.NewMySQL("")
	if err != nil {
		log.WithError(err).Error("failed to get mysql config")
		return
	}
	srv := migrate.DefaultMigrateService(db, r)
	ctx := life.SignalContext()
	for _, f := range flist {
		if err := f(srv, ctx); err != nil {
			log.Error(err)
			return
		}
	}
}
