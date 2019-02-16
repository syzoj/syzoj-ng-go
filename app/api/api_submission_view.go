package api

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"

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
		bson.D{{"_id", submissionId}, {"public", true}},
		mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"type", 1}, {"result.status", 1}, {"result.score", 1}, {"content.language", 1}, {"submit_time", 1}, {"problem", 1}, {"content.code", 1}}),
	).Decode(submission); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrSubmissionNotFound
		}
	}
	problem := new(model.Problem)
	if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(),
		bson.D{{"_id", submission.Problem}},
		mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"title", 1}}),
	).Decode(problem); err != nil {
		log.WithField("submissionId", submissionId).WithError(err).Warning("Failed to read problem for submission")
	}
	resp := new(model.SubmissionViewResponse)
	resp.Submission = submission
	resp.Problem = problem
	c.SendValue(resp)
	return nil
}
