package api

import (
	"github.com/valyala/fastjson"

    "github.com/syzoj/syzoj-ng-go/app/core"
)

// POST /api/problem-db/new
//
// Example request:
//     {
//         "problem": {
//             "title": "PRoblem Title"
//         }
//     }
//
// Example response:
//     {
//          "problem_id": "AAAAAAAAAAAAAAAA"
//     }
//
func Handle_ProblemDb_New(c *ApiContext) (apiErr ApiError) {
	var err error
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
	title := string(body.GetStringBytes("title"))
    resp, err := c.Server().c.Action_ProblemDb_New(c.Context(), &core.ProblemDbNew1{
        Title: title,
        Owner: c.Session.AuthUserUid,
    })
    switch err {
    case core.ErrInvalidProblem:
        return badRequestError(err)
    case nil:
        arena := new(fastjson.Arena)
        result := arena.NewObject()
        result.Set("problem_id", arena.NewString(EncodeObjectID(resp.ProblemId)))
        c.SendValue(result)
        return
    default:
        panic(err)
    }
}
