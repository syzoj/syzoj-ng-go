package judge

import (
	"errors"
	model_problem "github.com/syzoj/syzoj-ng-go/app/model/problem"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type JudgeService interface {
	CreateProblem(id util.UUID) error
	GetProblemStatement(id util.UUID) (model_problem.ProblemStatement, error)
	GetProblemPushToken(id util.UUID) (string, error)
}

type JudgeServiceProvider interface {
	GetJudgeService(name string) JudgeService
}

var ProblemNotFoundError = errors.New("Problem not found")
