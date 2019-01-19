package api

import (
	"time"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/google/uuid"
	"github.com/valyala/fastjson"
)

func Handle_ProblemDb_New(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	var body *fastjson.Value
	if body, err = c.GetBody(); err != nil {
		return badRequestError(err)
	}
	problemId := uuid.New()
	title := string(body.GetStringBytes("title"))
	var datetimeVal []byte
	if datetimeVal, err = time.Now().MarshalBinary(); err != nil {
		panic(err)
	}
	if _, err = c.Dgraph().NewTxn().Mutate(c.Context(), &dgo_api.Mutation{
		Set: []*dgo_api.NQuad{
			{
				Subject:     "_:problem",
				Predicate:   "problem.id",
				ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: problemId.String()}},
			},
			{
				Subject:   "_:problem",
				Predicate: "problem.owner",
				ObjectId:  c.Session.AuthUserUid,
			},
			{
				Subject:     "_:problem",
				Predicate:   "problem.title",
				ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: title}},
			},
			{
				Subject:     "_:problem",
				Predicate:   "problem.create_time",
				ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_DatetimeVal{DatetimeVal: datetimeVal}},
			},
		},
		CommitNow: true,
	}); err != nil {
		return internalServerError(err)
	}
	c.SendValue(fastjson.NewObject(map[string]*fastjson.Value{
		"problem_id": fastjson.NewString(problemId.String()),
	}))
	return
}
