package api

import (
	"time"

	"github.com/valyala/fastjson"
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/bson/primitive"
)

// POST /api/problem-db/new
//
// Example request:
//     {
//         "problem": {
//             "title": "PRoblem Title"
//         }
//     }
//
// Example response:
//     {
//          "problem_id": "AAAAAAAAAAAAAAAA"
//     }
//
func Handle_ProblemDb_New(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
        panic(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	var body *fastjson.Value
	if body, err = c.GetBody(); err != nil {
		return badRequestError(err)
	}
    problemId := primitive.NewObjectID()
    title := string(body.GetStringBytes("title"))
    if _, err = c.Server().mongodb.Collection("problem").InsertOne(c.Context(),
        bson.D{{"_id", problemId}, {"title", title}, {"owner", []primitive.ObjectID{c.Session.AuthUserUid}}, {"create_time", time.Now()}},
    ); err != nil {
        panic(err)
    }
    arena := new(fastjson.Arena)
    result := arena.NewObject()
    result.Set("problem_id", arena.NewString(EncodeObjectID(problemId)))
    c.SendValue(result)
	return
}
