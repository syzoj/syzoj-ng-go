package legacy

import (
	"errors"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

var ErrNoScore = errors.New("No score")

func (r *SubmissionResult) Model_GetScore() (float64, error) {
	if r.Score == nil {
		return 0, ErrNoScore
	}
	return r.GetScore(), nil
}

var _ model.SubmissionScore = &SubmissionResult{}
