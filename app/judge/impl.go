package judge

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"

	"github.com/dgraph-io/dgo"
	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"

	judge_api "github.com/syzoj/syzoj-ng-go/app/judge/protos"
)

var log = logrus.StandardLogger()

type judgeService struct {
	dgraph     *dgo.Dgraph
	queue      chan int64
	queueSize  int64
	queueItems sync.Map
}

type queueItem struct {
	Uid string `json:"uid"`
}

func (item *queueItem) getFields() logrus.Fields {
	return logrus.Fields{
		"Uid": item.Uid,
	}
}

func CreateJudgeService(dgraph *dgo.Dgraph) (Service, error) {
	srv := &judgeService{
		dgraph: dgraph,
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

	const query = `
{
	submission(func: has(submission.status)) {
		uid
	}
}
`
	var apiResponse *dgo_api.Response
	apiResponse, err = srv.dgraph.NewTxn().Query(context.TODO(), query)
	if err != nil {
		return
	}
	type response struct {
		Submission []*queueItem `json:"submission"`
	}
	var resp response
	if err = json.Unmarshal(apiResponse.Json, &resp); err != nil {
		return
	}
	log.WithField("submissionCount", len(resp.Submission)).Info("Adding stored submissions to queue")
	go func() {
		for _, v := range resp.Submission {
			srv.enqueue(v)
		}
	}()
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

func (srv *judgeService) NotifySubmission(ctx context.Context, uid string) (err error) {
	log.WithField("Uid", uid).Info("Received notify submission")
	const query = `
query SubmissionQuery($id: string) {
	submission(func: uid($id)) @filter(has(submission.status)) {
		uid
	}
}
`
	var apiResponse *dgo_api.Response
	apiResponse, err = srv.dgraph.NewTxn().QueryWithVars(context.TODO(), query, map[string]string{"$id": uid})
	if err != nil {
		return
	}
	type response struct {
		Submission []*queueItem `json:"submission"`
	}
	var resp response
	if err = json.Unmarshal(apiResponse.Json, &resp); err != nil {
		return
	}
	if len(resp.Submission) == 0 {
		panic("Invalid submission uid")
	}
	srv.enqueue(resp.Submission[0])
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
