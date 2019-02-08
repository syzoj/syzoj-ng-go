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
	lock       sync.RWMutex
	context    context.Context
	cancelFunc func()

	queue      chan int
	queueSize  int
	queueItems map[int]*queueItem
	queueLock  sync.Mutex

	// TODO: This may leak memory, perform GC on it
	submissions     map[primitive.ObjectID]*Submission
	submissionsLock sync.Mutex

	judgers    map[primitive.ObjectID]*judger
	judgerLock sync.Mutex

	contests map[primitive.ObjectID]*Contest
	wg       sync.WaitGroup

	oracle     map[interface{}]struct{}
	oracleLock sync.Mutex
}

func NewCore(mongodb *mongo.Client) (srv *Core, err error) {
	srv = &Core{
		mongodb: mongodb.Database("syzoj"),
	}
	srv.context, srv.cancelFunc = context.WithCancel(context.Background())
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if err = srv.initJudge(srv.context); err != nil {
		return
	}
	if err = srv.initContestLocked(srv.context); err != nil {
		return
	}
	srv.initOracle()
	return
}

func (c *Core) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cancelFunc()
	c.unloadAllContestsLocked()
	return nil
}
