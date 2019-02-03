package api

import (
	"strconv"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Contests(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	var cursor *mongo.Cursor
	if cursor, err = c.Server().mongodb.Collection("problemset").Find(c.Context(),
		bson.D{{"contest", bson.D{{"$exists", true}}}},
		mongo_options.Find().SetProjection(bson.D{{"_id", 1}, {"problemset_name", 1}, {"contest.start_time", 1}, {"contest.running", 1}}),
	); err != nil {
		panic(err)
	}
	defer cursor.Close(c.Context())
	arena := new(fastjson.Arena)
	result := arena.NewObject()
	contests := arena.NewArray()
	i := 0
	for cursor.Next(c.Context()) {
		var contest model.Problemset
		if err = cursor.Decode(&contest); err != nil {
			panic(err)
		}
		value := arena.NewObject()
		value.Set("id", arena.NewString(EncodeObjectID(contest.Id)))
		value.Set("title", arena.NewString(contest.ProblemsetName))
		if contest.Contest.Running {
			value.Set("running", arena.NewTrue())
		} else {
			value.Set("running", arena.NewFalse())
		}
		value.Set("start_time", arena.NewNumberString(strconv.FormatInt(contest.Contest.StartTime.Unix(), 10)))
		contests.SetArrayItem(i, value)
		i++
	}
	if err = cursor.Err(); err != nil {
		panic(err)
	}
	result.Set("contests", contests)
	c.SendValue(result)
	return
}
