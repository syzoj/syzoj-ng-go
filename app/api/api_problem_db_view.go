package api

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

// GET /api/problem-db/view/{problem_id}
//
// Path parameters:
//     problem_id: The ObjectID of the problem.
//
// Example response:
//     {
//         "title": "Problem Title",
//         "statement": "Problem Statement",
//         "is_owner": false,
//         "can_submit": false
//     }
func Handle_ProblemDb_View(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	problemId := model.MustDecodeObjectID(vars["problem_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	problem := new(model.Problem)
	if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(), bson.D{{"_id", problemId}}).Decode(problem); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrProblemNotFound
		}
		panic(err)
	}
	var isOwner bool
	if c.Session.LoggedIn() {
		for _, v := range problem.Owner {
			if model.MustGetObjectID(v) == c.Session.AuthUserUid {
				isOwner = true
				break
			}
		}
	}
	if !isOwner && !problem.GetPublic() {
		return ErrPermissionDenied
	}
	resp := new(model.ProblemDbNewResponse)
	resp.Problem = new(model.Problem)
	resp.Problem.Title = problem.Title
	resp.Problem.Statement = problem.Statement
	resp.CanSubmit = proto.Bool(c.Session.LoggedIn())
	resp.IsOwner = proto.Bool(isOwner)
	c.SendValue(resp)
	return nil
}

func Handle_ProblemDb_View_Edit(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	problemId := model.MustDecodeObjectID(vars["problem_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	body := new(model.ProblemDbEditRequest)
	if err = c.GetBody(body); err != nil {
		return badRequestError(err)
	}
	if body.Problem == nil {
		return ErrGeneral
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	problem := new(model.Problem)
	if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(), bson.D{{"_id", problemId}}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"owner", 1}})).Decode(problem); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrProblemNotFound
		}
		panic(err)
	}
	var isOwner bool
	for _, owner := range problem.Owner {
		if model.MustGetObjectID(owner) == c.Session.AuthUserUid {
			isOwner = true
		}
	}
	if !isOwner {
		return ErrPermissionDenied
	}
	newProblem := new(model.Problem)
	newProblem.Title = body.Problem.Title
	newProblem.Statement = body.Problem.Statement
	if newProblem.Title != nil || newProblem.Statement != nil {
		if _, err = c.Server().mongodb.Collection("problem").UpdateOne(c.Context(), bson.D{{"_id", problemId}}, bson.D{{"$set", newProblem}}); err != nil {
			panic(err)
		}
	}
	c.SendValue(&empty.Empty{})
	return nil
}

func Handle_ProblemDb_View_Submit(c *ApiContext) ApiError {
	var err error
	vars := c.Vars()
	problemId := model.MustDecodeObjectID(vars["problem_id"])
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	body := new(model.ProblemDbViewSubmitRequest)
	if err = c.GetBody(body); err != nil {
		return badRequestError(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	problem := new(model.Problem)
	if err = c.Server().mongodb.Collection("problem").FindOne(c.Context(),
		bson.D{{"_id", problemId}},
		mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"public", 1}}),
	).Decode(problem); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrProblemNotFound
		}
		panic(err)
	}
	if !problem.GetPublic() {
		return ErrPermissionDenied
	}

	submission := &model.Submission{
		Id:         model.NewObjectIDProto(),
		Type:       proto.String("standard"),
		User:       model.ObjectIDProto(c.Session.AuthUserUid),
		Owner:      []*model.ObjectID{model.ObjectIDProto(c.Session.AuthUserUid)},
		Problem:    model.ObjectIDProto(problemId),
		Content:    body.Code,
		SubmitTime: ptypes.TimestampNow(),
		Public:     proto.Bool(true),
	}
	if _, err = c.Server().mongodb.Collection("submission").InsertOne(c.Context(), submission); err != nil {
		panic(err)
	}
	go c.Server().c.EnqueueSubmission(model.MustGetObjectID(submission.Id))
	go c.Server().c.IncrementSubmissionCounter(problemId, model.MustGetObjectID(submission.Id))
	resp := new(model.ProblemDbViewSubmitResponse)
	resp.Submission = new(model.Submission)
	resp.Submission.Id = submission.Id
	c.SendValue(resp)
	return nil
}
