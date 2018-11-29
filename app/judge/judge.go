package judge

import (
	"errors"

	"github.com/google/uuid"
	model_problem "github.com/syzoj/syzoj-ng-go/app/model/problem"
)

type JudgeService interface {
	CreateProblem(id uuid.UUID) error
	GetProblemStatement(id uuid.UUID) (model_problem.ProblemStatement, error)
	GetProblemPushToken(id uuid.UUID) (string, error)
}

type JudgeServiceProvider interface {
	GetJudgeService(name string) JudgeService
}

var ProblemNotFoundError = errors.New("Problem not found")
