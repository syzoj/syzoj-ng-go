package main

import (
	"context"
	"fmt"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/config"
	"github.com/syzoj/syzoj-ng-go/svc/core"
	"github.com/syzoj/syzoj-ng-go/svc/redis"
	"github.com/volatiletech/sqlboiler/boil"
)

var log = logrus.StandardLogger()

func main() {
	boil.DebugMode = true
	ctx := context.Background()
	db, err := config.NewMySQL("")
	if err != nil {
		panic(err)
	}

	r, err := config.NewRedis("")
	if err != nil {
		panic(err)
	}
	rs := redis.DefaultRedisService(r)
	var tmr *redis.RedisTimer
	tmr = rs.DefaultTimer("timer", func(key string) {
		log.Info(key)
		tmr.Delete(ctx, key)
	})
	//tmr.Schedule(ctx, "hello2", time.Now().Add(time.Second*2))
	log.Info("start")
	tmr.Run(ctx)
	return

	c := core.DefaultCore(db, rs)
	uid, v, err := c.CreateProblem(ctx, map[string]interface{}{
		"hello": "world5",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(uid, v)
	v2, err := c.UpdateProblemCAS(ctx, uid, map[string]interface{}{
		"hello": "world6",
	}, 1)
	if err != nil {
		panic(err)
	}
	meta, v3, err := c.GetProblemMetadata(ctx, uid+"x")
	if err != nil {
		panic(err)
	}
	if v2 != v3 {
		panic("version updated")
	}
	fmt.Println(uid, string(meta), v3)
}
