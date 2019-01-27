package api

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/valyala/fastjson"
)

// POST /api/nav/logout
//
// Does not take a body.
//
// If logout succeeds, returns `nil`. Otherwise, returns an error indicating reason for failure/
func Handle_Nav_Logout(c *ApiContext) (apiErr ApiError) {
	var err error
	if err = c.SessionStart(); err != nil {
		panic(err)
	}
	if !c.Session.LoggedIn() {
		return ErrNotLoggedIn
	}
	if _, err = c.Server().mongodb.Collection("session").UpdateOne(c.Context(),
		bson.D{{"_id", c.Session.SessUid}},
		bson.D{{"$unset", bson.D{{"session_user", nil}}}},
	); err != nil {
		panic(err)
	}
	if err = c.SessionReload(); err != nil {
		panic(err)
	}
	c.SendValue(new(fastjson.Arena).NewNull())
	return
}
