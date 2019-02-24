package core

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (c *Core) IncrementSubmissionCounter(problemId primitive.ObjectID, submissionId primitive.ObjectID) {
	_, err := c.mongodb.Collection("problem").UpdateOne(c.context, bson.D{{"_id", problemId}}, bson.D{{"$inc", bson.D{{"public_stats.submission", 1}}}})
	if err != nil {
		log.WithField("problemId", problemId).WithField("submissionId", submissionId).WithError(err).Error("Failed to increment submission count")
	}
}
