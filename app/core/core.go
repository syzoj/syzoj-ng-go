package core

import (
	"context"
	"sync"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

type Core struct {
	mongodb    *mongo.Database
	lock       sync.Mutex
	context    context.Context
	cancelFunc func()

	queue      chan int64
	queueSize  int64
	queueItems sync.Map

	contests map[primitive.ObjectID]*Contest
	wg       sync.WaitGroup
}

func NewCore(mongodb *mongo.Client) (srv *Core, err error) {
	srv = &Core{
		mongodb: mongodb.Database("syzoj"),
	}
	srv.context, srv.cancelFunc = context.WithCancel(context.Background())
	if err = srv.initJudge(srv.context); err != nil {
		return
	}
	if err = srv.initContest(srv.context); err != nil {
		return
	}
	return
}

func (srv *Core) Close() error {
	srv.cancelFunc()
	srv.unloadAllContests()
	return nil
}
