package api

import (
	"strconv"
	"sync"

	"github.com/mongodb/mongo-go-driver/bson"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Contest_Submissions(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	contestId := DecodeObjectID(vars["contest_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	contest := c.Server().c.GetContestR(contestId)
	if contest == nil {
		return ErrContestNotLoaded
	}
	player := contest.GetPlayer(c.Session.AuthUserUid)
	if player == nil {
		contest.RUnlock()
		return ErrGeneral // user is not player
	}
	arena := new(fastjson.Arena)
	result := arena.NewObject()
	var wg sync.WaitGroup
	func() {
		defer contest.RUnlock()
		playerProblems := player.GetProblems()
		problemsValue := arena.NewArray()
		problemI := 0
		for name, problem := range playerProblems {
			problemValue := arena.NewObject()
			problemValue.Set("name", arena.NewString(name))
			submissionsValue := arena.NewArray()
			for submissionI, submission := range problem.GetSubmissions() {
				submissionValue := arena.NewObject()
				submissionValue.Set("id", arena.NewString(EncodeObjectID(submission.GetSubmissionId())))
				submissionValue.Set("submit_time", arena.NewNumberString(strconv.FormatInt(submission.GetSubmitTime().Unix(), 10)))
				submissionsValue.SetArrayItem(submissionI, submissionValue)
			}
			contestProblem := contest.GetProblemByName(name)
			problemValue.Set("submissions", submissionsValue)
			wg.Add(1)
			go func() {
				defer wg.Done()
				var problemModel model.Problem
				var arena fastjson.Arena
				if err := c.Server().mongodb.Collection("problem").FindOne(c.Context(), bson.D{{"_id", contestProblem.ProblemId}}, mongo_options.FindOne().SetProjection(bson.D{{"title", 1}})).Decode(&problemModel); err != nil {
					log.WithField("problemId", contestProblem.ProblemId).Error("Failed to read contest problem title: ", err)
					problemValue.Set("title", arena.NewString("<ERROR>"))
				} else {
					problemValue.Set("title", arena.NewString(problemModel.Title))
				}
			}()
			problemsValue.SetArrayItem(problemI, problemValue)
			problemI++
		}
		result.Set("problems", problemsValue)
	}()
	wg.Wait()
	c.SendValue(result)
	return nil
}
