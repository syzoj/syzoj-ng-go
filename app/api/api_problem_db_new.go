package api

import (
	"github.com/golang/protobuf/ptypes"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

func Handle_ProblemDb_New(c *ApiContext) ApiError {
	var err error
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	body := new(model.ProblemDbNewRequest)
	if err = c.GetBody(body); err != nil {
		return badRequestError(err)
	}
	newProblem := new(model.Problem)
	newProblem.Id = model.NewObjectIDProto()
	newProblem.Owner = []*model.ObjectID{model.ObjectIDProto(c.Session.AuthUserUid)}
	newProblem.Title = body.Title
	newProblem.Statement = body.Statement
	newProblem.CreateTime = ptypes.TimestampNow()
	if _, err = c.Server().mongodb.Collection("problem").InsertOne(c.Context(), newProblem); err != nil {
		panic(err)
	}
	resp := new(model.ProblemDbNewResponse)
	resp.ProblemId = newProblem.Id
	c.SendValue(resp)
	return nil
}
