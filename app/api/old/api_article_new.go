package api

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/valyala/fastjson"
)

// POST /api/article/new
//
// Example request:
//     {
//         "article": {
//             "title": "Title",
//             "text": "Text"
//         }
//     }
//
// Example response:
//     {
//         "id": "AAAAAAAAAAAAAAAA"
//     }
func Handle_Article_New(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		return
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	var body *fastjson.Value
	if body, err = c.GetBody(); err != nil {
		return badRequestError(err)
	}
	articleId := primitive.NewObjectID()
	title := string(body.GetStringBytes("article", "title"))
	text := string(body.GetStringBytes("article", "text"))
	if _, err = c.Server().mongodb.Collection("article").InsertOne(c.Context(), bson.D{
		{"_id", articleId},
		{"title", title},
		{"owner", c.Session.AuthUserUid},
		{"text", text},
		{"reply", bson.A{}},
		{"create_time", time.Now()},
		{"last_edit_time", time.Now()},
	}); err != nil {
		panic(err)
	}
	arena := new(fastjson.Arena)
	result := arena.NewObject()
	result.Set("id", arena.NewString(EncodeObjectID(articleId)))
	c.SendValue(result)
	return
}
