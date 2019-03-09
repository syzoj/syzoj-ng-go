package judger

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"

	"github.com/syzoj/syzoj-ng-go/judger/rpc"
)

type judgeContextKey struct{}

type JudgeContext struct {
	j   *Judger
	ctx context.Context
	tag int64
}

func GetJudgeContext(ctx context.Context) *JudgeContext {
	c := ctx.Value(judgeContextKey{})
	if c == nil {
		return nil
	}
	return c.(*JudgeContext)
}

func (c *JudgeContext) ReportResult(result *any.Any) error {
	_, err := c.j.client.SetTaskResult(c.ctx, &rpc.SetTaskResultMessage{
		JudgerId:    c.j.judgeRequest.JudgerId,
		JudgerToken: c.j.judgeRequest.JudgerToken,
		TaskTag:     proto.Int64(c.tag),
		Result:      result,
	})
	return err
}

func (c *JudgeContext) ReportProgress(result *any.Any) error {
	return nil
}
