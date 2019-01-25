package api

import (
	"time"

	"github.com/google/uuid"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/valyala/fastjson"
	"golang.org/x/crypto/bcrypt"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

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
	var user model.User
	user.Id = primitive.NewObjectID()
	user.UserName = &userName
	user.RegisterTime = time.Now()
	user.Auth = new(model.UserAuth)
	xid := uuid.New()
	user.Xid = &xid
	if user.Auth.Password, err = bcrypt.GenerateFromPassword([]byte(password), 0); err != nil {
		panic(err)
	}
	if _, err = c.Server().mongodb.Collection("user").InsertOne(c.Context(), user); err != nil {
		panic(err)
	}
	log.WithField("username", userName).Info("Created account")
	c.SendValue(new(fastjson.Arena).NewNull())
	return
}
