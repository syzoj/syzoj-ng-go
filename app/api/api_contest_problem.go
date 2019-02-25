package api

import (
    "go.mongodb.org/mongo-driver/bson"

    "github.com/syzoj/syzoj-ng-go/app/model"
    "github.com/syzoj/syzoj-ng-go/app/core"
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
