package api

import (
    "github.com/mongodb/mongo-go-driver/bson"
    "github.com/mongodb/mongo-go-driver/mongo"
    mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
    "github.com/valyala/fastjson"

    "github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Submissions(c *ApiContext) (apiErr ApiError) {
    var err error
    if err = c.SessionStart(); err != nil {
        panic(err)
    }

    query := bson.D{}
    if c.FormValue("my") != "" {
        if c.Session.LoggedIn() {
            query = append(query, bson.E{"user", c.Session.AuthUserUid})
        }
    }
    var cursor mongo.Cursor
    if cursor, err = c.Server().mongodb.Collection("submission").Find(c.Context(), query,
        mongo_options.Find().SetProjection(bson.D{{"_id", 1}, {"problem", 1}, {"result.status", 1}, {"result.score", 1}, {"user", 1}, {"content.language", 1}, {"submit_time", 1}}).SetLimit(50).SetSort(bson.D{{"submit_time", -1}})); err != nil {
        panic(err)
    }
    defer cursor.Close(c.Context())

    arena := new(fastjson.Arena)
    result := arena.NewObject()
    submissions := arena.NewArray()
    item := 0
    for cursor.Next(c.Context()) {
        var submission model.Submission
        if err = cursor.Decode(&submission); err != nil {
            return
        }
        value := arena.NewObject()
        value.Set("id", arena.NewString(EncodeObjectID(submission.Id)))
        value.Set("problem_id", arena.NewString(EncodeObjectID(submission.Problem)))
        value.Set("problem_title", arena.NewString("TODO"))
        value.Set("status", arena.NewString(submission.Result.Status))
        value.Set("score", arena.NewNumberFloat64(submission.Result.Score))
        value.Set("language", arena.NewString(submission.Content.Language))
        value.Set("submit_user_id", arena.NewString(EncodeObjectID(submission.User)))
        value.Set("submit_user_name", arena.NewString("TODO"))
        value.Set("submit_time", arena.NewString(submission.SubmitTime.String()))
        submissions.SetArrayItem(item, value)
        item++
    }
    result.Set("submissions", submissions)
    c.SendValue(result)
    return
}
