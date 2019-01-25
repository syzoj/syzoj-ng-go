package judge

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"

	judge_api "github.com/syzoj/syzoj-ng-go/app/judge/protos"
)

var log = logrus.StandardLogger()

type judgeService struct {
	mongodb    *mongo.Client
	queue      chan int64
	queueSize  int64
	queueItems sync.Map
}

type queueItem struct {
}

func (item *queueItem) getFields() logrus.Fields {
	return nil
}

func CreateJudgeService(mongodb *mongo.Client) (Service, error) {
	srv := &judgeService{
		mongodb: mongodb,
	}
	err := srv.init()
	if err != nil {
		return srv, err
	}
	return srv, nil
}

func (srv *judgeService) init() (err error) {
	srv.queue = make(chan int64, 1000)
	srv.queueItems = sync.Map{}
	srv.queueSize = 0
	// TODO: init queue
	return
}

func (srv *judgeService) Close() error {
	return nil
}

func (srv *judgeService) enqueue(item *queueItem) {
	log.WithFields(item.getFields()).Info("Adding submission to queue")
	i := atomic.AddInt64(&srv.queueSize, 1)
	srv.queueItems.Store(i, item)
	srv.queue <- i
}

func (srv *judgeService) NotifySubmission(ctx context.Context, id string) (err error) {
	log.WithField("id", id).Info("Received notify submission, TODO: enqueue")
	return
}

func (srv *judgeService) RegisterJudger(req *judge_api.JudgeRequest, s judge_api.Judge_RegisterJudgerServer) (err error) {
	return
}

func (srv *judgeService) SetTaskProgress(s judge_api.Judge_SetTaskProgressServer) (err error) {
	return
}

func (srv *judgeService) SetTaskResult(ctx context.Context, in *judge_api.TaskResult) (e *empty.Empty, err error) {
	return
}
