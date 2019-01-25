package api

import (
    mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
    "github.com/mongodb/mongo-go-driver/mongo"
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/valyala/fastjson"

    "github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Submission_View(c *ApiContext) (apiErr ApiError) {
    var err error
    vars := c.Vars()
    submissionId := DecodeObjectID(vars["submission_id"])
    if err = c.SessionStart(); err != nil {
        return
    }
    var submissionModel model.Submission
    if err = c.Server().mongodb.Collection("submission").FindOne(c.Context(),
        bson.D{{"_id", submissionId}},
        mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"type", 1}, {"status", 1}, {"score", 1}, {"language", 1}, {"submit_time", 1}, {"problem", 1}, {"code", 1}}),
    ).Decode(&submissionModel); err != nil {
        if err == mongo.ErrNoDocuments {
            return ErrSubmissionNotFound
        }
    }
    log.Info(submissionModel)
    var problemModel model.Problem
    var problemTitle string
    if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(),
        bson.D{{"_id", submissionModel.Problem}},
        mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"title", 1}}),
    ).Decode(&problemModel); err != nil {
    } else {
        problemTitle = problemModel.Title
    }
    arena := new(fastjson.Arena)
    result := arena.NewObject()
    submission := arena.NewObject()
    submission.Set("status", arena.NewString(submissionModel.Status))
    submission.Set("score", arena.NewNumberFloat64(submissionModel.Score))
    submission.Set("language", arena.NewString(submissionModel.Language))
    submission.Set("problem_id", arena.NewString(EncodeObjectID(submissionModel.Problem)))
    submission.Set("problem_title", arena.NewString(problemTitle))
    submission.Set("submit_time", arena.NewString(submissionModel.SubmitTime.String()))
    submission.Set("code", arena.NewString(submissionModel.Code))
    result.Set("submission", submission)
    c.SendValue(result)
    return
}
