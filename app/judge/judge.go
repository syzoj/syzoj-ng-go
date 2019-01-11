package judge

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	judge_api "github.com/syzoj/syzoj-ng-go/app/judge/protos"
)

// The interface for judge queue.
type Service interface {
	// Notify the service of a new submission.
	NotifySubmission(ctx context.Context, uid string) error
	Close() error

	RegisterJudger(*judge_api.JudgeRequest, judge_api.Judge_RegisterJudgerServer) error
	SetTaskProgress(judge_api.Judge_SetTaskProgressServer) error
	SetTaskResult(ctx context.Context, in *judge_api.TaskResult) (*empty.Empty, error)
}
