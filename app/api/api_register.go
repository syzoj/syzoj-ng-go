package api

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/syzoj/syzoj-ng-go/app/core"
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
	body := new(model.RegisterRequest)
	if err = c.GetBody(body); err != nil {
		return badRequestError(err)
	}

	userName := body.GetUsername()
	if !model.CheckName(userName) {
		return ErrInvalidUserName
	}
	password := body.GetPassword()
	var passwordHash []byte
	if passwordHash, err = bcrypt.GenerateFromPassword([]byte(password), 0); err != nil {
		panic(err)
	}
	email := body.GetEmail()
	if !model.CheckEmail(email) {
		return ErrInvalidEmail
	}
	lock := c.Server().c.LockOracle([]interface{}{core.KeyUserName(userName), core.KeyEmail(email)})
	if lock == nil {
		return ErrConflict
	}
	defer lock.Release()
	var n int64
	if n, err = c.Server().mongodb.Collection("user").CountDocuments(c.Context(), bson.D{{"username", userName}}, mongo_options.Count().SetLimit(1)); err != nil {
		panic(err)
	} else if n != 0 {
		return ErrDuplicateUserName
	}
	if n, err = c.Server().mongodb.Collection("user").CountDocuments(c.Context(), bson.D{{"email", email}}, mongo_options.Count().SetLimit(1)); err != nil {
		panic(err)
	} else if n != 0 {
		return ErrDuplicateEmail
	}
	userModel := &model.User{
		Id:           model.NewObjectIDProto(),
		Username:     proto.String(userName),
		Email:        proto.String(email),
		RegisterTime: ptypes.TimestampNow(),
		Auth: &model.UserAuth{
			Method:   proto.Int64(1),
			Password: passwordHash,
		},
	}
	if _, err = c.Server().mongodb.Collection("user").InsertOne(c.Context(), userModel); err != nil {
		panic(err)
	}
	log.WithField("username", userName).Info("Created account")
	c.SendValue(&empty.Empty{})
	return
}
