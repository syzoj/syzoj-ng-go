package api

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"
	"golang.org/x/crypto/bcrypt"

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
	if !checkUserName(userName) {
		return ErrInvalidUserName
	}
	password := string(body.GetStringBytes("password"))
	var passwordHash []byte
	if passwordHash, err = bcrypt.GenerateFromPassword([]byte(password), 0); err != nil {
		panic(err)
	}
	lock := c.Server().c.LockOracle([]interface{}{core.KeyUserName(userName)})
	if lock == nil {
		return ErrConflict
	}
	defer lock.Release()
	if _, err = c.Server().mongodb.Collection("user").FindOne(c.Context(), bson.D{{"username", userName}}, mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}})).DecodeBytes(); err != nil {
		if err != mongo.ErrNoDocuments {
			panic(err)
		}
	} else {
		return ErrDuplicateUserName
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
