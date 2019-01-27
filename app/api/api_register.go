package api

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/valyala/fastjson"
	"golang.org/x/crypto/bcrypt"
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
	if !checkUserName(userName) {
		return ErrInvalidUserName
	}
	password := string(body.GetStringBytes("password"))
    var passwordHash []byte
	if passwordHash, err = bcrypt.GenerateFromPassword([]byte(password), 0); err != nil {
		panic(err)
	}
    userId := primitive.NewObjectID()
	if _, err = c.Server().mongodb.Collection("user").InsertOne(c.Context(), bson.D{
        {"_id", userId},
        {"username", userName},
        {"register_time", time.Now()},
        {"auth", bson.D{{"password", passwordHash}, {"method", int64(1)}}},
    }); err != nil {
		panic(err)
	}
	log.WithField("username", userName).Info("Created account")
	c.SendValue(new(fastjson.Arena).NewNull())
	return
}
