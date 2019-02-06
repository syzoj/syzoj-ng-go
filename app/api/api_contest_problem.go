package api

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/core"
	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Contest_Problem(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	contestId := DecodeObjectID(vars["contest_id"])
	entryName := vars["entry_name"]
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	contest := c.Server().c.GetContestR(contestId)
	if contest == nil {
		return ErrContestNotFound
	}
	running := contest.Running()
	if !running {
		return ErrPermissionDenied
	}
	entryId, found := contest.NameToProblems[entryName]
	if !found {
		return ErrContestNotFound
	}
	problemsetEntry := contest.Problems[entryId]
	contest.RUnlock()

	var problemModel model.Problem
	if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(), bson.D{{"_id", problemsetEntry.ProblemId}}, mongo_options.FindOne().SetProjection(bson.D{{"title", 1}, {"statement", 1}})).Decode(&problemModel); err != nil {
		if err == mongo.ErrNoDocuments {
			log.WithField("contestId", contestId).WithField("entryName", entryName).WithField("problemId", problemsetEntry.ProblemId).Error("Problem referenced by contest does not exist")
		}
		panic(err)
	}

	arena := new(fastjson.Arena)
	result := arena.NewObject()
	problem := arena.NewObject()
	problem.Set("title", arena.NewString(problemModel.Title))
	problem.Set("statement", arena.NewString(problemModel.Statement))
	result.Set("problem", problem)
	c.SendValue(result)
	return
}

func Handle_Contest_Problem_Submit(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	contestId := DecodeObjectID(vars["contest_id"])
	entryName := vars["entry_name"]
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	var body *fastjson.Value
	body, err = c.GetBody()
	if err != nil {
		return badRequestError(err)
	}

	contest := c.Server().c.GetContestR(contestId)
	if contest == nil {
		return ErrContestNotFound
	}
	running := contest.Running()
	if !running {
		contest.RUnlock()
		return ErrPermissionDenied
	}
	player := contest.GetPlayer(c.Session.AuthUserUid)
	if player == nil {
		contest.RUnlock()
		return ErrPermissionDenied
	}
	entryId, found := contest.NameToProblems[entryName]
	if !found {
		contest.RUnlock()
		return ErrContestNotFound
	}
	problemsetEntry := contest.Problems[entryId]
	contest.RUnlock()

	var resp *core.Submit1Resp
	switch resp, err = c.Server().c.Action_Submit(c.Context(), &core.Submit1{
		ProblemId: problemsetEntry.ProblemId,
		Submitter: c.Session.AuthUserUid,
		Enqueue:   false,
		Public:    false,
		Code: core.Code{
			Language: string(body.GetStringBytes("code", "language")),
			Code:     string(body.GetStringBytes("code", "code")),
		},
	}); err {
	case core.ErrProblemNotExist:
		return ErrContestNotFound
	case nil:
	default:
		panic(err)
	}

	contest = c.Server().c.GetContestW(contestId)
	if contest == nil {
		return ErrContestNotFound
	}
	running = contest.Running()
	if !running {
		return ErrPermissionDenied
	}
	player = contest.GetPlayer(c.Session.AuthUserUid)
	if player == nil {
		return ErrPermissionDenied
	}
	err = contest.PlayerSubmission(player, entryName, resp.SubmissionId)
	contest.Unlock()
	switch err {
	case core.ErrGeneral:
		return internalServerError(err)
	case core.ErrTooManySubmissions:
		return ErrTooManySubmissions
	}
	arena := new(fastjson.Arena)
	result := arena.NewObject()
	result.Set("submission_id", arena.NewString(EncodeObjectID(resp.SubmissionId)))
	c.SendValue(result)
	return nil
}
