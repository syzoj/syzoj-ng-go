package core

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type contestHook struct {
	*Contest
}

func (ct contestHook) OnSubmissionResult(submissionId primitive.ObjectID, result *any.Any) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.submissionResults[submissionId] = result
	ct.notifyUpdateRanklist()
}

func (ct *Contest) AppendSubmission(player *ContestPlayer, name string, submissionId primitive.ObjectID) {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	var problem *ContestPlayerProblem
	{
		var found bool
		for _, problem = range player.problems {
			if problem.name == name {
				found = true
				break
			}
		}
		if !found {
			problem = new(ContestPlayerProblem)
			problem.name = name
			player.problems = append(player.problems, problem)
		}
	}
	submission := &ContestPlayerProblemSubmission{submissionId: submissionId}
	problem.submissions = append(problem.submissions, submission)
	// note the race condition here
	go ct.fetchSubmission(submissionId)
	if ct.judgeInContest {
		go ct.c.EnqueueSubmission(submissionId)
	}
}

func (ct *Contest) fetchSubmission(submissionId primitive.ObjectID) {
	submission := new(model.Submission)
	if err := ct.c.mongodb.Collection("submission").FindOne(ct.c.context, bson.D{{"_id", submissionId}}).Decode(submission); err != nil {
		log.WithField("submissionId", submissionId).WithError(err).Error("Failed to fetch submission result")
		return
	}
	var dany ptypes.DynamicAny
	if submission.Result != nil {
		if err := ptypes.UnmarshalAny(submission.Result, &dany); err != nil {
			log.WithField("submissionId", submissionId).WithError(err).Error("Failed to parse submission result")
		}
	}
	ct.mu.Lock()
	defer ct.mu.Unlock()
	ct.submissionResults[submissionId] = dany.Message
	ct.notifyUpdateRanklist()
}

func (ct *Contest) GetRanklist() *model.ContestRanklist {
	ct.mu.Lock()
	defer ct.mu.Unlock()
	return ct.ranklist
}
