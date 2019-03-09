package judger

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/sirupsen/logrus"
	"github.com/syzoj/syzoj-ng-go/judger/rpc"
)

var log = logrus.StandardLogger()

type Judger struct {
	rpcUrl       string
	judgeRequest *rpc.JudgeRequest

	client rpc.JudgeClient
}

func NewJudger() *Judger {
	return new(Judger)
}

func (j *Judger) Run(ctx context.Context) error {
	conn, err := grpc.Dial(j.rpcUrl)
	if err != nil {
		return err
	}
	defer conn.Close()
	j.client = rpc.NewJudgeClient(conn)

	var timeout time.Duration
	for {
		task, err := j.client.FetchTask(ctx, j.judgeRequest)
		if err != nil {
			log.WithError(err).Warningf("Failed to fetch task, trying after %d seconds")
			time.Sleep(timeout)
			timeout = timeout * 2
			if timeout > time.Minute {
				timeout = time.Minute
			}
			continue
		} else {
			timeout = time.Millisecond * 10
		}
		if task.GetSuccess() {
			j.handleTask(ctx, task.GetTask())
		}
	}
	return nil
}

func (j *Judger) handleTask(ctx context.Context, task *rpc.Task) {
	typ := task.GetProblemType()
	backend := backends[typ]
	if backend == nil {
		log.WithField("task_tag", task.GetTaskTag()).WithField("backend_type", typ).Error("Unrecognized backend type")
		return
	}
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	c := &JudgeContext{
		j:   j,
		ctx: ctx2,
		tag: task.GetTaskTag(),
	}
	err := backend.JudgeSubmission(context.WithValue(ctx2, judgeContextKey{}, c), task)
	if err != nil {
		log.WithField("task_tag", task.GetTaskTag()).WithField("backend_type", typ).WithError(err).Error("Judgement failed")
		return
	}
}
