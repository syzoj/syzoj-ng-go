package judge

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
    "github.com/mongodb/mongo-go-driver/bson/primitive"

	judge_api "github.com/syzoj/syzoj-ng-go/app/judge/protos"
)

type SubmissionId = primitive.ObjectID

// The interface for judge queue.
type Service interface {
	// Notify the service of a new submission.
	NotifySubmission(ctx context.Context, id SubmissionId) error
	Close() error

    FetchTask(context.Context, *judge_api.JudgeRequest) (*judge_api.FetchTaskResult, error)
	SetTaskProgress(judge_api.Judge_SetTaskProgressServer) error
	SetTaskResult(ctx context.Context, in *judge_api.TaskResult) (*empty.Empty, error)
}
