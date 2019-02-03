package api

import (
	"github.com/mongodb/mongo-go-driver/bson"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"

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
	var contestModel model.Problemset
	if err = c.Server().mongodb.Collection("problemset").FindOne(c.Context(), bson.D{
		{"_id", contestId},
		{"contest", bson.D{{"$exists", true}}},
	}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"description", 1}})).Decode(&contestModel); err != nil {
		log.Info("not found in mongodb")
		return ErrContestNotFound
	}

	contest := c.Server().c.GetContestR(contestId)
	if contest == nil {
		log.Info("not found in memory")
		return ErrContestNotFound
	}
	running := contest.Running()
	contest.RUnlock()

	arena := new(fastjson.Arena)
	result := arena.NewObject()
	contestObj := arena.NewObject()
	contestObj.Set("description", arena.NewString(contestModel.Description))
	if running {
		contestObj.Set("running", arena.NewTrue())
	} else {
		contestObj.Set("running", arena.NewFalse())
	}
	result.Set("contest", contestObj)
	c.SendValue(result)
	return
}
