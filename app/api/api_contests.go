package api

import (
	"strconv"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Contests(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	var cursor *mongo.Cursor
	if cursor, err = c.Server().mongodb.Collection("contest").Find(c.Context(),
		bson.D{},
	); err != nil {
		panic(err)
	}
	defer cursor.Close(c.Context())
	arena := new(fastjson.Arena)
	result := arena.NewObject()
	contests := arena.NewArray()
	i := 0
	for cursor.Next(c.Context()) {
		var contestModel model.Contest
		if err = cursor.Decode(&contestModel); err != nil {
			panic(err)
		}
		value := arena.NewObject()
		value.Set("id", arena.NewString(EncodeObjectID(contestModel.Id)))
		value.Set("title", arena.NewString(contestModel.Name))
		if contestModel.Running {
			value.Set("running", arena.NewTrue())
		} else {
			value.Set("running", arena.NewFalse())
		}
		value.Set("start_time", arena.NewNumberString(strconv.FormatInt(contestModel.StartTime.Unix(), 10)))
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
