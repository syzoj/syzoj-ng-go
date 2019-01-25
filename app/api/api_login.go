package api

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	mongo_options "github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/valyala/fastjson"
	"golang.org/x/crypto/bcrypt"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

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
	if user.Auth == nil {
		return ErrCannotLogin
	}
	if err = bcrypt.CompareHashAndPassword(user.Auth.Password, []byte(password)); err != nil {
		return ErrPasswordIncorrect
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
