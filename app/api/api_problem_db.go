package api

import (
	"regexp"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

// GET /api/problem-db
//
// Query parameters:
//     my: if exists, show only problems by myself (requires login)
//     search: if exists, a substring to match in problem title
//     limit: An integer, the max number of problems to return, max 100
//     skip: An integer, how many documents to skip, default 0
//
// Response: A `problems` array with each object corresponding to a problem in the results, and the count of matching documents.
//
// Example response:
//     {
//         "problems": [
//             {
//                 "id": "AAAAAAAAAAAAAAAA",
//                 "title": "Problem Title",
//                 "create_time": 0,
//             }
//          ]
//      }
//
//
func Handle_ProblemDb(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	query := bson.D{}
	form := c.Form()
	if len(form["my"]) != 0 {
		if !c.Session.LoggedIn() {
			return ErrNotLoggedIn
		}
		query = append(query, bson.E{"owner", c.Session.AuthUserUid})
	} else {
		query = append(query, bson.E{"public", true})
	}
	if len(form["search"]) != 0 {
		query = append(query, bson.E{"title", bson.D{{"$regex", regexp.QuoteMeta(form["search"][0])}}})
	}
	var cursor *mongo.Cursor
	options := mongo_options.Find()
	options.SetProjection(bson.D{{"_id", "1"}, {"title", 1}, {"create_time", 1}, {"public_stats.submission", 1}})
	if len(form["skip"]) != 0 {
		skip, err := strconv.ParseInt(form["skip"][0], 10, 64)
		if err == nil {
			options.SetSkip(skip)
		}
	}
	var limit int64 = 100
	if len(form["limit"]) != 0 {
		l, err := strconv.ParseInt(form["limit"][0], 10, 64)
		if err == nil {
			limit = l
		}
	}
	if limit > 100 {
		limit = 100
	}
	options.SetLimit(limit)
	if cursor, err = c.Server().mongodb.Collection("problem").Find(c.Context(), query, options); err != nil {
		panic(err)
	}
	defer cursor.Close(c.Context())

	resp := new(model.ProblemDbResponse)
	for cursor.Next(c.Context()) {
		problem := new(model.Problem)
		if err = cursor.Decode(problem); err != nil {
			return
		}
		resp.Problems = append(resp.Problems, problem)
	}
	if err = cursor.Err(); err != nil {
		panic(err)
	}
	c.SendValue(resp)
	return
}
