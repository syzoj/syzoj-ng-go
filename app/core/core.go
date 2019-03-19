package core

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var log = logrus.StandardLogger()

type Core struct {
	mongodb *mongo.Database
	lock    sync.RWMutex
	// context is the reliability context, aborting it would result in data loss
	context context.Context
	// context2 is thoe cancel signal context, aborting it should only cause minimal damages
	context2    context.Context
	cancelFunc  func()
	cancelFunc2 func()

	queue                chan int
	queueSize            int
	queueItems           map[int]*queueItem
	queueLock            sync.Mutex
	judgers              map[primitive.ObjectID]*judger
	judgerLock           sync.Mutex
	submissionHooks      map[SubmissionHook]struct{}
	submissionHooksMutex sync.Mutex

	wg sync.WaitGroup

	oracle     map[interface{}]struct{}
	oracleLock sync.Mutex
}

func NewCore(mongodb *mongo.Client) (srv *Core, err error) {
	srv = &Core{
		mongodb: mongodb.Database("syzoj"),
	}
	srv.submissionHooks = make(map[SubmissionHook]struct{})
	srv.context, srv.cancelFunc = context.WithCancel(context.Background())
	srv.context2, srv.cancelFunc2 = context.WithCancel(context.Background())
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if err = srv.initJudge(srv.context); err != nil {
		return
	}
	srv.initOracle()
	return
}

func (c *Core) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cancelFunc2()
	return nil
}
