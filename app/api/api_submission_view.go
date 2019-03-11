package api

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_Submission_View(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	submissionId := model.MustDecodeObjectID(vars["submission_id"])
	if err = c.SessionStart(); err != nil {
		return
	}
	submission := new(model.Submission)
	if err = c.Server().mongodb.Collection("submission").FindOne(c.Context(),
		bson.D{{"_id", submissionId}, {"public", true}}).Decode(submission); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrSubmissionNotFound
		}
		panic(err)
	}
	var inf interface{}
	if err = c.Server().mongodb.Collection("submission").FindOne(c.Context(),
		bson.D{{"_id", submissionId}, {"public", true}}).Decode(&inf); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrSubmissionNotFound
		}
		panic(err)
	}
	problem := new(model.Problem)
	if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(),
		bson.D{{"_id", submission.Problem}},
		mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"title", 1}}),
	).Decode(problem); err != nil {
		log.WithField("submissionId", submissionId).WithError(err).Warning("Failed to read problem for submission")
	}
	resp := new(model.SubmissionViewResponse)
	resp.SubmissionContent = submission.Content
	resp.SubmissionResult = submission.Result
	resp.ProblemId = problem.Id
	resp.ProblemTitle = problem.Title
	c.SendValue(resp)
	return nil
}
