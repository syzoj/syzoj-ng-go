package main

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	_ "github.com/syzoj/syzoj-ng-go/judger/backend/legacy"
	"github.com/syzoj/syzoj-ng-go/judger/judger"
)

var log = logrus.StandardLogger()

func main() {
	app := cli.NewApp()
	app.Name = "syzoj-ng-judge"
	app.Usage = "syzoj-ng judge daemon"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: "config.yaml",
			Usage: "The config file for judger",
		},
	}
	app.Action = func(c *cli.Context) error {
		j := judger.NewJudger()
		if err := j.LoadConfig(c.String("config")); err != nil {
			log.Fatal(err)
		}
		if err := j.Run(context.Background()); err != nil {
			log.Fatal(err)
		}
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
