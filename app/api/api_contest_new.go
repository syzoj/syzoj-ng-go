package api

import (
	"errors"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/core"
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
	problemsetId := primitive.NewObjectID()
	title := string(body.GetStringBytes("title"))
	if _, err = c.Server().mongodb.Collection("problemset").InsertOne(c.Context(), bson.D{{"_id", problemsetId}, {"problemset_name", title}}); err != nil {
		panic(err)
	}
	var options core.ContestOptions
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
	if err = c.Server().c.CreateContest(c.Context(), problemsetId, &options); err != nil {
		if err == core.ErrInvalidOptions {
			return badRequestError(err)
		} else {
			panic(err)
		}
	}
	arena := new(fastjson.Arena)
	result := arena.NewObject()
	result.Set("id", arena.NewString(EncodeObjectID(problemsetId)))
	c.SendValue(result)
	return
}