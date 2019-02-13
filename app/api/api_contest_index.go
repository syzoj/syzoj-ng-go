package api

import (
	"sync"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/core"
	"github.com/syzoj/syzoj-ng-go/app/model"
)

// GET /api/contest/{contest_id}/index
func Handle_Contest_Index(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	contestId := DecodeObjectID(vars["contest_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	var contestModel model.Contest
	if err = c.Server().mongodb.Collection("contest").FindOne(c.Context(), bson.D{
		{"_id", contestId},
	}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"name", 1}, {"description", 1}})).Decode(&contestModel); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrContestNotFound
		}
		panic(err)
	}

	contest := c.Server().c.GetContestR(contestId)
	if contest == nil {
		return ErrContestNotLoaded
	}
	running := contest.Running()
	var player *core.ContestPlayer
	if c.Session.LoggedIn() {
		player = contest.GetPlayer(c.Session.AuthUserUid)
	}
	var problemsVisible bool
	contestProblems := contest.GetProblems()
	if running && player != nil {
		problemsVisible = true
	}
	contest.RUnlock()

	arena := new(fastjson.Arena)
	result := arena.NewObject()
	result.Set("name", arena.NewString(contestModel.Name))
	result.Set("description", arena.NewString(contestModel.Description))
	contestObj := arena.NewObject()
	if running {
		contestObj.Set("running", arena.NewTrue())
	} else {
		contestObj.Set("running", arena.NewFalse())
	}
	if problemsVisible {
		problems := arena.NewArray()
		var wg sync.WaitGroup
		problemsModel := make([]model.Problem, len(contestProblems))
		for i, problemEntry := range contestProblems {
			wg.Add(1)
			go func(i int, id primitive.ObjectID) {
				defer wg.Done()
				var err error
				if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(), bson.D{{"_id", id}}, mongo_options.FindOne().SetProjection(bson.D{{"title", 1}})).Decode(&problemsModel[i]); err != nil {
					log.Error("Failed to query problem in subquery: ", err)
					return
				}
			}(i, problemEntry.ProblemId)
		}
		wg.Wait()
		for i := range contestProblems {
			problem := arena.NewObject()
			problem.Set("problem_id", arena.NewString(EncodeObjectID(problemsModel[i].Id)))
			problem.Set("entry_name", arena.NewString(contestProblems[i].Name))
			problem.Set("title", arena.NewString(problemsModel[i].Title))
			problems.SetArrayItem(i, problem)
		}
		contestObj.Set("problems", problems)
	} else {
		contestObj.Set("running", arena.NewFalse())
	}
	result.Set("contest", contestObj)
	c.SendValue(result)
	return
}
