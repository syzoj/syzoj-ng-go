package api

import (
	"sync"

	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/syzoj/syzoj-ng-go/app/core"
	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Contest_Index(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	contestId := model.MustDecodeObjectID(vars["contest_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	contestModel := new(model.Contest)
	if err = c.Server().mongodb.Collection("contest").FindOne(c.Context(), bson.D{{"_id", contestId}}).Decode(contestModel); err != nil {
		panic(err)
	}
	contest := c.Server().c.GetContest(contestId)
	if contest == nil {
		return ErrContestNotFound
	}
	problems := contest.GetProblems()
	resp := new(model.ContestIndexResponse)
	resp.Contest = new(model.Contest)
	resp.Contest.Name = contestModel.Name
	resp.Contest.Description = contestModel.Description
	resp.Running = proto.Bool(contest.IsRunning())
	var player *core.ContestPlayer
	if c.Session.LoggedIn() {
		player = contest.GetPlayerById(c.Session.AuthUserUid)
	}
	if contest.CheckListProblems(player) {
		var wg sync.WaitGroup
		resp.Problems = make([]*model.ContestProblemEntryResponse, len(problems))
		for i, problem := range problems {
			entry := new(model.ContestProblemEntryResponse)
			entry.Name = proto.String(problem.GetName())
			wg.Add(1)
			go func(problemId *model.ObjectID) {
				defer wg.Done()
				problemModel := new(model.Problem)
				if err := c.Server().mongodb.Collection("problem").FindOne(c.Context(), bson.D{{"_id", problemId}}).Decode(problemModel); err != nil {
					log.WithField("problemId", problemId).WithError(err).Error("Failed to get problem title")
					return
				}
				entry.Problem = new(model.Problem)
				entry.Problem.Title = problemModel.Title
			}(problem.GetData().Problem)
			resp.Problems[i] = entry
		}
		wg.Wait()
	}
	c.SendValue(resp)
	return nil
}
