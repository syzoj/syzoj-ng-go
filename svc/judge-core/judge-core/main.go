package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/lib/config"
	"github.com/syzoj/syzoj-ng-go/svc/judge-core"
)

var log = logrus.StandardLogger()

func main() {
	port, err := config.GetHttpListenPort()
	if err != nil {
		log.Fatalf("failed to get http listen port: %s", err)
	}
	minio, err := config.NewMinio("")
	if err != nil {
		log.Fatalf("failed to create minio client: %s", err)
	}
	s := &judgecore.Server{
		ListenAddr: fmt.Sprintf(":%d", port),
		Minio:      minio,
		FileBucket: "judge",
	}
	if err := s.Serve(); err != nil {
		log.Error(err)
	}
}
