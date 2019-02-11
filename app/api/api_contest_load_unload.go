package api

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Contest_Load(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	contestId := DecodeObjectID(vars["contest_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	var contestModel model.Contest
	if err = c.Server().mongodb.Collection("contest").FindOne(c.Context(), bson.D{{"_id", contestId}}, mongo_options.FindOne().SetProjection(bson.D{{"owner", 1}})).Decode(&contestModel); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrContestNotFound
		}
		panic(err)
	}
	if contestModel.Owner != c.Session.AuthUserUid {
		return ErrPermissionDenied
	}
	if err = c.Server().c.LoadContest(contestId); err != nil {
		panic(err)
	}
	arena := new(fastjson.Arena)
	c.SendValue(arena.NewNull())
	return nil
}

func Handle_Contest_Unload(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	contestId := DecodeObjectID(vars["contest_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	var contestModel model.Contest
	if err = c.Server().mongodb.Collection("contest").FindOne(c.Context(), bson.D{{"_id", contestId}}, mongo_options.FindOne().SetProjection(bson.D{{"owner", 1}})).Decode(&contestModel); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrContestNotFound
		}
		panic(err)
	}
	if contestModel.Owner != c.Session.AuthUserUid {
		return ErrPermissionDenied
	}
	if err = c.Server().c.UnloadContest(contestId); err != nil {
		panic(err)
	}
	arena := new(fastjson.Arena)
	c.SendValue(arena.NewNull())
	return nil
}
