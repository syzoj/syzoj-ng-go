package api

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_P(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	shortName := vars["short_name"]
	var problemModel model.Problem
	if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(), bson.D{{"short_name", shortName}}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}})).Decode(&problemModel); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrProblemNotFound
		}
	}
	arena := new(fastjson.Arena)
	result := arena.NewObject()
	result.Set("problem_id", arena.NewString(EncodeObjectID(problemModel.Id)))
	c.SendValue(result)
	return nil
}
