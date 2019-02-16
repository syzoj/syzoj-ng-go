package api

import (
	"sync"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Submissions(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		panic(err)
	}

	query := bson.D{{"public", true}}
	form := c.Form()
	if len(form["my"]) != 0 {
		if c.Session.LoggedIn() {
			query = append(query, bson.E{"user", c.Session.AuthUserUid})
		}
	}
	if len(form["problem"]) != 0 {
		var ids []primitive.ObjectID
		for _, problemIdStr := range form["problem"] {
			id, err := model.DecodeObjectID(problemIdStr)
			if err != nil {
				ids = append(ids, id)
			}
		}
		if len(ids) != 0 {
			query = append(query, bson.E{"problem", bson.D{{"$in", ids}}})
		}
	}
	var cursor *mongo.Cursor
	if cursor, err = c.Server().mongodb.Collection("submission").Find(c.Context(), query,
		mongo_options.Find().SetProjection(bson.D{{"problem", 1}, {"result.status", 1}, {"result.score", 1}, {"user", 1}, {"content.language", 1}, {"submit_time", 1}}).SetLimit(50).SetSort(bson.D{{"submit_time", -1}})); err != nil {
		panic(err)
	}
	defer cursor.Close(c.Context())

	resp := new(model.SubmissionsResponse)
	var wg sync.WaitGroup
	for cursor.Next(c.Context()) {
		submission := new(model.Submission)
		if err = cursor.Decode(submission); err != nil {
			return
		}
		entry := &model.SubmissionsResponseSubmissionEntry{
			Submission: submission,
			Problem:    &model.Problem{},
			SubmitUser: &model.User{},
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := c.Server().mongodb.Collection("problem").FindOne(c.Context(), bson.D{{"_id", submission.Problem}}, mongo_options.FindOne().SetProjection(bson.D{{"title", 1}})).Decode(entry.Problem)
			if err != nil {
				log.WithField("submissionId", submission.Id).WithError(err).Warning("Failed to get problem for submission")
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := c.Server().mongodb.Collection("user").FindOne(c.Context(), bson.D{{"_id", submission.User}}, mongo_options.FindOne().SetProjection(bson.D{{"username", 1}})).Decode(entry.SubmitUser)
			if err != nil {
				log.WithField("submissionId", submission.Id).WithError(err).Warning("Failed to get user for submission")
			}
		}()
		resp.Submissions = append(resp.Submissions, entry)
	}
	wg.Wait()
	c.SendValue(resp)
	return
}
