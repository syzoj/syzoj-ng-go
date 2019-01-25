package api

import (
    "github.com/mongodb/mongo-go-driver/mongo"
    mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/valyala/fastjson"

    "github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_ProblemDb_View(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
    problemId := DecodeObjectID(vars["problem_id"])
	if err = c.SessionStart(); err != nil {
        panic(err)
	}
    var problem model.Problem
    if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(), bson.D{{"_id", problemId}}).Decode(&problem); err != nil {
        if err == mongo.ErrNoDocuments {
            return ErrProblemNotFound
        }
        panic(err)
    }
    arena := new(fastjson.Arena)
    result := arena.NewObject()
    result.Set("title", arena.NewString(problem.Title))
    result.Set("statement", arena.NewString(problem.Statement))
    if c.Session.LoggedIn() {
        result.Set("can_submit", arena.NewTrue())
    } else {
        result.Set("can_submit", arena.NewFalse())
    }
    var isOwner bool
	if c.Session.LoggedIn() {
        for _, v := range problem.Owner {
            if v == c.Session.AuthUserUid {
                isOwner = true
                break
            }
        }
    }
    if isOwner {
        result.Set("is_owner", arena.NewTrue())
	} else {
        result.Set("is_owner", arena.NewFalse())
    }
	c.SendValue(result)
	return
}

func Handle_ProblemDb_View_Submit(c *ApiContext) (apiErr ApiError) {
    var err error
    vars := c.Vars();
    problemId := DecodeObjectID(vars["problem_id"])
    var problemModel model.Problem
    if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(),
        bson.D{{"_id", problemId}},
        mongo_options.FindOne().SetProjection(bson.D{{"x", "y"}}),
    ).Decode(&problemModel); err != nil {
        if err == mongo.ErrNoDocuments {
            return ErrProblemNotFound
        }
        panic(err)
    }
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
    submissionId := primitive.NewObjectID()
    if _, err = c.Server().mongodb.Collection("submission").InsertOne(c.Context(), bson.D{
        {"_id", submissionId},
        {"type", "standard"},
        {"user", c.Session.AuthUserUid},
        {"owner", []primitive.ObjectID{c.Session.AuthUserUid}},
        {"problem", problemId},
        {"language", string(body.GetStringBytes("code", "language"))},
        {"code", string(body.GetStringBytes("code", "code"))},
    }); err != nil {
        return
    }
    arena := new(fastjson.Arena)
    result := arena.NewObject()
    result.Set("id", arena.NewString(EncodeObjectID(submissionId)))
    c.SendValue(result);
    return
}
