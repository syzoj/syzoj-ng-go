package api

import (
	"time"

	dgo_api "github.com/dgraph-io/dgo/protos/api"
	"github.com/google/uuid"
	"github.com/valyala/fastjson"
)

func Handle_ProblemDb_View(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	var problemId = uuid.MustParse(vars["problem_id"])
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	var dgValue *fastjson.Value
	var userId string
	if c.Session.LoggedIn() {
		userId = c.Session.AuthUserUid
	} else {
		userId = "0x0"
	}
	if dgValue, err = c.Query(ViewProblemQuery, map[string]string{"$problemId": problemId.String(), "$userId": userId}); err != nil {
		return internalServerError(err)
	}
	if len(dgValue.GetArray("problem")) == 0 {
		return ErrProblemNotFound
	}
	problemValue := dgValue.Get("problem", "0")
	resp := make(map[string]*fastjson.Value)
	resp["title"] = fastjson.NewString(string(problemValue.GetStringBytes("title")))
	resp["statement"] = fastjson.NewString(string(problemValue.GetStringBytes("statement")))
	resp["can_submit"] = fastjson.NewBool(c.Session.LoggedIn())
	resp["can_publicize"] = fastjson.NewBool(len(dgValue.GetArray("problemset")) > 0)
	if c.Session.LoggedIn() && string(problemValue.GetStringBytes("owner_uid")) == c.Session.AuthUserUid {
		resp["is_owner"] = fastjson.NewBool(true)
		resp["token"] = fastjson.NewString(string(problemValue.GetStringBytes("token")))
	} else {
		resp["is_owner"] = fastjson.NewBool(false)
	}
	c.SendValue(fastjson.NewObject(resp))
	return
}

func Handle_ProblemDb_View_Submit(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	problemId := uuid.MustParse(vars["problem_id"])
	var body *fastjson.Value
	if body, err = c.GetBody(); err != nil {
		return badRequestError(err)
	}
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	var dgValue *fastjson.Value
	if dgValue, err = c.Query(CheckProblemCanSubmitQuery, map[string]string{"$problemId": problemId.String()}); err != nil {
		return
	}
	if len(dgValue.GetArray("problem")) == 0 {
		return ErrProblemNotFound
	}
	submissionId := uuid.New()
	var datetimeVal []byte
	if datetimeVal, err = time.Now().MarshalBinary(); err != nil {
		return
	}
	var assigned *dgo_api.Assigned
	if assigned, err = c.Dgraph().NewTxn().Mutate(c.Context(), &dgo_api.Mutation{
		Set: []*dgo_api.NQuad{
			{
				Subject:     "_:submission",
				Predicate:   "submission.id",
				ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: submissionId.String()}},
			},
			{
				Subject:   "_:submission",
				Predicate: "submission.problem",
				ObjectId:  string(dgValue.GetStringBytes("problem", "0", "uid")),
			},
			{
				Subject:     "_:submission",
				Predicate:   "submission.language",
				ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: string(body.GetStringBytes("language"))}},
			},
			{
				Subject:     "_:submission",
				Predicate:   "submission.code",
				ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: string(body.GetStringBytes("code"))}},
			},
			{
				Subject:   "_:submission",
				Predicate: "submission.owner",
				ObjectId:  c.Session.AuthUserUid,
			},
			{
				Subject:     "_:submission",
				Predicate:   "submission.submit_time",
				ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_DatetimeVal{DatetimeVal: datetimeVal}},
			},
			{
				Subject:     "_:submission",
				Predicate:   "submission.status",
				ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: "Waiting"}},
			},
		},
		CommitNow: true,
	}); err != nil {
		return internalServerError(err)
	}
	c.SendValue(fastjson.NewObject(map[string]*fastjson.Value{
		"submission_id": fastjson.NewString(submissionId.String()),
	}))
	if err := c.JudgeService().NotifySubmission(c.Context(), assigned.Uids["submission"]); err != nil {
		log.Error(err)
	}
	return
}

func Handle_ProblemDb_View_Publicize(c *ApiContext) (apiErr ApiError) {
	var err error
	vars := c.Vars()
	var problemId = uuid.MustParse(vars["problem_id"])
	if err = c.SessionStart(); err != nil {
		return internalServerError(err)
	}
	var body *fastjson.Value
	if body, err = c.GetBody(); err != nil {
		return badRequestError(err)
	}
	name := string(body.GetStringBytes("name"))
	if name == "" {
		return badRequestError(ErrInvalidPublicName)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	if err = c.DgraphTransaction(func(t *DgraphTransaction) (err error) {
		var dgResponse *dgo_api.Response
		if dgResponse, err = c.Dgraph().NewReadOnlyTxn().QueryWithVars(c.Context(), CheckProblemPublicizeQuery, map[string]string{"$problemId": problemId.String(), "$userId": c.Session.AuthUserUid, "$name": name}); err != nil {
			return err
		}
		dgParser := c.GetParser()
		defer c.PutParser(dgParser)
		var dgValue *fastjson.Value
		if dgValue, err = dgParser.ParseBytes(dgResponse.Json); err != nil {
			panic(err)
		}
		if len(dgValue.GetArray("problem")) == 0 {
			t.Defer(func() {
				apiErr = ErrProblemNotFound
				return
			})
			return
		}
		if string(dgValue.GetStringBytes("problem", "0", "owner_uid")) != c.Session.AuthUserUid || len(dgValue.GetArray("problemset")) == 0 {
			t.Defer(func() {
				apiErr = ErrPermissionDenied
				return
			})
			return
		}
		if len(dgValue.GetArray("name")) != 0 {
			t.Defer(func() {
				apiErr = ErrDuplicatePublicName
				return
			})
			return
		}
		_, err = t.T.Mutate(c.Context(), &dgo_api.Mutation{
			Set: []*dgo_api.NQuad{
				{
					Subject:     "_:problemsetentry",
					Predicate:   "problemsetentry.id",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: uuid.New().String()}},
				},
				{
					Subject:     "_:problemsetentry",
					Predicate:   "problemsetentry.short_name",
					ObjectValue: &dgo_api.Value{Val: &dgo_api.Value_StrVal{StrVal: name}},
				},
				{
					Subject:   "_:problemsetentry",
					Predicate: "problemsetentry.problem",
					ObjectId:  string(dgValue.GetStringBytes("problem", "0", "uid")),
				},
				{
					Subject:   "_:problemsetentry",
					Predicate: "problemsetentry.problemset",
					ObjectId:  string(dgValue.GetStringBytes("problemset", "0", "uid")),
				},
			},
		})
		if err != nil {
			return
		}
		t.Defer(func() {
			c.SendValue(dgParser.NewNull())
		})
		return
	}); err != nil {
		return internalServerError(err)
	}
	return
}
