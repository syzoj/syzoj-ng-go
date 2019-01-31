package api

import (
	"strconv"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

// GET /api/articles
func Handle_Articles(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		return
	}
	var cursor mongo.Cursor
	if cursor, err = c.Server().mongodb.Collection("article").Find(c.Context(),
		bson.D{},
		mongo_options.Find().SetProjection(bson.D{
			{"_id", 1},
			{"title", 1},
			{"owner", 1},
			{"create_time", 1},
		}).SetLimit(50)); err != nil {
		panic(err)
	}
	arena := new(fastjson.Arena)
	result := arena.NewObject()
	articles := arena.NewArray()
	item := 0
	for cursor.Next(c.Context()) {
		var articleModel model.Article
		if err = cursor.Decode(&articleModel); err != nil {
			panic(err)
		}
		article := arena.NewObject()
		article.Set("id", arena.NewString(EncodeObjectID(articleModel.Id)))
		article.Set("title", arena.NewString(articleModel.Title))
		article.Set("owner_id", arena.NewString(EncodeObjectID(articleModel.Owner)))
		article.Set("create_time", arena.NewNumberString(strconv.FormatInt(articleModel.CreateTime.Unix(), 10)))
		articles.SetArrayItem(item, article)
		item++
	}
	result.Set("articles", articles)
	if err = cursor.Err(); err != nil {
		panic(err)
	}
	c.SendValue(result)
	return
}
