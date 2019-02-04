package api

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"

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
	var contestModel model.Problemset
	if err = c.Server().mongodb.Collection("problemset").FindOne(c.Context(), bson.D{
		{"_id", contestId},
		{"contest", bson.D{{"$exists", true}}},
	}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"problems", 1}})).Decode(&contestModel); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrContestNotFound
		}
		panic(err)
	}
	contest := c.Server().c.GetContestR(contestId)
	if contest == nil {
		return ErrContestNotFound
	}
	running := contest.Running()
	contest.RUnlock()
	if !running {
		return ErrPermissionDenied
	}

	var problemsetEntry model.ProblemsetEntry
	var found bool
	for _, problemsetEntry = range contestModel.Problems {
		if problemsetEntry.Name == entryName {
			found = true
			break
		}
	}
	if !found {
		return ErrContestNotFound
	}

	var problemModel model.Problem
	if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(), bson.D{{"_id", problemsetEntry.ProblemId}}, mongo_options.FindOne().SetProjection(bson.D{{"title", 1}, {"statement", 1}})).Decode(&problemModel); err != nil {
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
