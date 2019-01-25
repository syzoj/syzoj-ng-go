package api

import (
    "errors"
    "time"

    "github.com/mongodb/mongo-go-driver/bson/primitive"
    "github.com/valyala/fastjson"
    
    "github.com/syzoj/syzoj-ng-go/app/model"
    "github.com/syzoj/syzoj-ng-go/app/contest"
)

func Handle_Contest_New(c *ApiContext) (apiErr ApiError) {
    var err error
    if err = c.SessionStart(); err != nil {
        panic(err)
    }
    var body *fastjson.Value
    if body, err = c.GetBody(); err != nil {
        return badRequestError(err)
    }
    var m model.Problemset
    m.Id = primitive.NewObjectID()
    title := string(body.GetStringBytes("title"))
    m.ProblemsetName = &title
    if _, err = c.Server().mongodb.Collection("problemset").InsertOne(c.Context(), m); err != nil {
        panic(err)
    }
    var options contest.ContestOptions
    optionsVal := body.Get("options")
    if optionsVal == nil {
        return badRequestError(errors.New("Invalid options"))
    }
    log.Info(string(optionsVal.GetStringBytes("start_time")))
    if options.StartTime, err = time.Parse(time.RFC3339, string(optionsVal.GetStringBytes("start_time"))); err != nil {
        return badRequestError(errors.New("Invalid start time: " + err.Error()))
    }
    if options.Duration, err = time.ParseDuration(string(optionsVal.GetStringBytes("duration"))); err != nil {
        return badRequestError(errors.New("Invalid duration"))
    }
    options.Rules.JudgeInContest = optionsVal.GetBool("rules", "judge_in_contest")
    options.Rules.SeeResult = optionsVal.GetBool("rules", "see_result")
    options.Rules.RejudgeAfterContest = optionsVal.GetBool("rules", "rejudge_after_contest")
    options.Rules.RanklistType = string(optionsVal.GetStringBytes("rules", "ranklist_type"))
    options.Rules.RanklistVisibility = string(optionsVal.GetStringBytes("rules", "ranklist_visibility"))
    if err = c.Server().contestService.CreateContest(c.Context(), m.Id, &options); err != nil {
        if err == contest.ErrInvalidOptions {
            return badRequestError(err)
        } else {
            panic(err)
        }
    }
    arena := new(fastjson.Arena)
    result := arena.NewObject()
    result.Set("id", arena.NewString(EncodeObjectID(m.Id)))
    c.SendValue(result)
    return
}
