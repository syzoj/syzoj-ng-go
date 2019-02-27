package api

import (
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/syzoj/syzoj-ng-go/app/core"
	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Contest_Problem_View(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	contestId := model.MustDecodeObjectID(vars["contest_id"])
	problemName := vars["problem_name"]
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	contest := c.Server().c.GetContest(contestId)
	if contest == nil {
		return ErrContestNotFound
	}
	var player *core.ContestPlayer
	if c.Session.LoggedIn() {
		player = contest.GetPlayerById(c.Session.AuthUserUid)
	}
	if !contest.CheckViewProblem(player, problemName) {
		return ErrPermissionDenied
	}
	problem := contest.GetProblemByName(problemName)
	if problem == nil {
		return ErrGeneral
	}
	resp := new(model.ContestProblemViewResponse)
	problemId, _ := model.GetObjectID(problem.GetData().Problem)
	problemModel := new(model.Problem)
	if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(), bson.D{{"_id", problemId}}).Decode(problemModel); err != nil {
		panic(err)
	}
	resp.Problem = new(model.Problem)
	resp.Problem.Title = problemModel.Title
	resp.Problem.Statement = problemModel.Statement
	c.SendValue(resp)
	return nil
}

func Handle_Contest_Problem_Submit(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	contestId := model.MustDecodeObjectID(vars["contest_id"])
	problemName := vars["problem_name"]
	req := new(model.ContestProblemSubmitRequest)
	if err = c.GetBody(req); err != nil {
		return badRequestError(err)
	}
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	contest := c.Server().c.GetContest(contestId)
	if contest == nil {
		return ErrContestNotFound
	}
	var player *core.ContestPlayer
	if c.Session.LoggedIn() {
		player = contest.GetPlayerById(c.Session.AuthUserUid)
	}
	if !contest.CheckSubmitProblem(player, problemName) {
		return ErrPermissionDenied
	}
	problem := contest.GetProblemByName(problemName)
	if problem == nil {
		return ErrGeneral
	}
	submissionModel := new(model.Submission)
	submissionModel.Id = model.NewObjectIDProto()
	submissionModel.Problem = problem.GetData().GetProblem()
	submissionModel.SubmitTime = ptypes.TimestampNow()
	submissionModel.Content = req.Content
	if _, err = c.Server().mongodb.Collection("submission").InsertOne(c.Context(), submissionModel); err != nil {
		panic(err)
	}
	contest.AppendSubmission(player, problemName, model.MustGetObjectID(submissionModel.Id))
	resp := new(model.ContestProblemSubmitResponse)
	resp.Submission = new(model.Submission)
	resp.Submission.Id = submissionModel.Id
	c.SendValue(resp)
	return nil
}
