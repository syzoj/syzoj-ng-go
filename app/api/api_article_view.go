package api

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

// GET /api/article/view/{article_id}
//
// Example response:
//     {
//         "article": {
//             "id": "AAAAAAAAAAAAAAAA",
//             "title": "Title",
//             "owner": "TODO",
//             "owner_id", "AAAAAAAAAAAAAAAA",
//             "create_time": "2006-01-02 15:04:05.999999999 -0700 MST"
//             "last_edit_time": "2006-01-02 15:04:05.999999999 -0700 MST",
//             "text": "Text"
//         }
//     }
func Handle_Article_View(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	articleId := DecodeObjectID(vars["article_id"])
	if err = c.SessionStart(); err != nil {
		return
	}
	var articleModel model.Article
	if err = c.Server().mongodb.Collection("article").FindOne(c.Context(),
		bson.D{{"_id", articleId}},
		mongo_options.FindOne().SetProjection(bson.D{
			{"_id", 1},
			{"title", 1},
			{"owner", 1},
			{"create_time", 1},
			{"last_edit_time", 1},
			{"text", 1},
			{"reply", 1},
		})).Decode(&articleModel); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNotFound
		}
		panic(err)
	}
	arena := new(fastjson.Arena)
	result := arena.NewObject()
	article := arena.NewObject()
	article.Set("id", arena.NewString(EncodeObjectID(articleModel.Id)))
	article.Set("title", arena.NewString(articleModel.Title))
	article.Set("owner", arena.NewString("TODO"))
	article.Set("owner_id", arena.NewString(EncodeObjectID(articleModel.Owner)))
	article.Set("create_time", arena.NewString(articleModel.CreateTime.String()))
	article.Set("last_edit_time", arena.NewString(articleModel.LastEditTime.String()))
	article.Set("text", arena.NewString(articleModel.Text))
	replies := arena.NewArray()
	itemReply := 0
	for _, replyModel := range articleModel.Reply {
		reply := arena.NewObject()
		reply.Set("owner", arena.NewString("TODO"))
		reply.Set("owner_id", arena.NewString(EncodeObjectID(replyModel.Owner)))
		reply.Set("text", arena.NewString(replyModel.Text))
		reply.Set("create_time", arena.NewString(replyModel.CreateTime.String()))
		replies.SetArrayItem(itemReply, reply)
	}
	article.Set("replies", replies)
	result.Set("article", article)
	c.SendValue(result)
	return
}
