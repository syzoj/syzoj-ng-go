package api

import (
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/core"
)

// POST /api/register
//
// Example request:
//     {
//         "username": "username",
//         "password": "password"
//     }
// If register succeeds, returns `nil`. Otherwise, returns an error indicating the reason for failure.
func Handle_Register(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	if c.Session.LoggedIn() {
		return ErrAlreadyLoggedIn
	}
	var body *fastjson.Value
	if body, err = c.GetBody(); err != nil {
		return badRequestError(err)
	}
	userName := string(body.GetStringBytes("username"))
	password := string(body.GetStringBytes("password"))
    _, err = c.Server().c.Action_Register(c.Context(), &core.Register1{
        UserName: userName,
        Password: password,
    })
    switch err {
    case core.ErrInvalidUserName:
        return ErrInvalidUserName
    case core.ErrDuplicateUserName:
        return ErrDuplicateUserName
    case nil:
	    c.SendValue(new(fastjson.Arena).NewNull())
    	return
    default:
        panic(err)
    }
}
