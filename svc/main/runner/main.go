package main

import (
	docker "github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/config"
	"github.com/syzoj/syzoj-ng-go/lib/life"
	lredis "github.com/syzoj/syzoj-ng-go/lib/redis"
	"github.com/syzoj/syzoj-ng-go/svc/judge/judge"
	"github.com/syzoj/syzoj-ng-go/svc/runner/runner"
)

var log = logrus.StandardLogger()

func main() {
	redisp, err := config.NewRedis("")
	if err != nil {
		log.Errorf("failed to create redis: %s", err)
		return
	}
	cli, err := docker.NewClientWithOpts(docker.FromEnv)
	if err != nil {
		log.Errorf("failed to create docker: %s", err)
		return
	}
	r := &runner.Runner{
		WorkingDir:        "/opt/runner/work",
		MountedWorkingDir: "/work",
		DataDir:           "/opt/runner/data",
		Docker:            cli,
		DockerImage:       "local/ubuntu",
		JudgeService:      judge.DefaultJudgeService(lredis.WrapPool(redisp)),
	}
	ctx := life.SignalContext()
	if err := r.Run(ctx); err != nil {
		log.Error(err)
	}
}
