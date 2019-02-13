package api

import (
	"crypto/md5"
	"crypto/subtle"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"
	"golang.org/x/crypto/bcrypt"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

// POST /api/login
//
// Example request:
//     {
//         "username": "username",
//         "password": "password"
//     }
// If login succeeds, returns `nil`. Otherwise, returns an error indicating the reason for failure.
func Handle_Login(c *ApiContext) (apiErr ApiError) {
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
	var user model.User
	if err = c.Server().mongodb.Collection("user").FindOne(c.Context(),
		bson.D{{"username", userName}},
		mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"auth", 1}}),
	).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrUserNotFound
		}
		panic(err)
	}
	switch user.Auth.Method {
	case 1:
		if err = bcrypt.CompareHashAndPassword(user.Auth.Password, []byte(password)); err != nil {
			return ErrPasswordIncorrect
		}
	case 2:
		sum := md5.Sum([]byte(password + "syzoj2_xxx"))
		if subtle.ConstantTimeCompare(sum[:], user.Auth.Password) != 1 {
			return ErrPasswordIncorrect
		}
	default:
		return ErrCannotLogin
	}
	if _, err = c.Server().mongodb.Collection("session").UpdateOne(c.Context(),
		bson.D{{"_id", c.Session.SessUid}},
		bson.D{{"$set", bson.D{{"session_user", user.Id}}}},
	); err != nil {
		panic(err)
	}
	c.SendValue(new(fastjson.Arena).NewNull())
	return
}
