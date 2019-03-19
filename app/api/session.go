package api

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	mongo_options "go.mongodb.org/mongo-driver/mongo/options"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type Session struct {
	SessUid          primitive.ObjectID
	HasAuthUserUid   bool
	AuthUserUid      primitive.ObjectID
	AuthUserUserName string
}

func (s *Session) LoggedIn() bool {
	return s.HasAuthUserUid
}

func (c *ApiContext) newSession(ctx context.Context) (err error) {
	var tokenBytes [16]byte
	rand.Read(tokenBytes[:])
	newToken := hex.EncodeToString(tokenBytes[:])
	session := new(model.Session)
	session.Id = model.NewObjectIDProto()
	session.SessionToken = proto.String(newToken)
	if _, err = c.Server.mongodb.Collection("session").InsertOne(ctx, session); err != nil {
		panic(err)
	}
	c.Session = new(Session)
	c.Session.SessUid = model.MustGetObjectID(session.Id)
	c.SetCookie(&http.Cookie{
		Name:     "SYZOJSESSION",
		HttpOnly: true,
		Path:     "/",
		Value:    newToken,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
	})
	return
}

func (c *ApiContext) SessionStart(ctx context.Context) (err error) {
	var claimedToken = c.GetCookie("SYZOJSESSION")
	if len(claimedToken) != 32 {
		claimedToken = ""
	}
	session := new(model.Session)
	if err = c.Server.mongodb.Collection("session").FindOne(ctx,
		bson.D{{"session_token", claimedToken}},
		mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"session_user", 1}}),
	).Decode(session); err == mongo.ErrNoDocuments {
		return c.newSession(ctx)
	} else if err != nil {
		panic(err)
	}
	c.Session = new(Session)
	c.Session.SessUid = model.MustGetObjectID(session.GetId())
	if session.SessionUser != nil {
		c.Session.HasAuthUserUid = true
		c.Session.AuthUserUid = model.MustGetObjectID(session.SessionUser)
		user := new(model.User)
		if err = c.Server.mongodb.Collection("user").FindOne(ctx,
			bson.D{{"_id", session.SessionUser}},
			mongo_options.FindOne().SetProjection(bson.D{{"_id", "1"}, {"username", 1}}),
		).Decode(user); err != nil {
			if err == mongo.ErrNoDocuments {
				log.WithField("SessUid", c.Session.SessUid).WithField("UserID", c.Session.AuthUserUid).Warning("Broken session: user doesn't exist")
			} else {
				panic(err)
			}
		}
		c.Session.AuthUserUserName = user.GetUsername()
	}
	return nil
}

func (c *ApiContext) SessionReload(ctx context.Context) (err error) {
	if c.Session == nil {
		panic("Calling SessionReload() without existing session")
	}
	session := new(model.Session)
	if err = c.Server.mongodb.Collection("session").FindOne(ctx,
		bson.D{{"_id", c.Session.SessUid}},
		mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"session_user", 1}}),
	).Decode(session); err == mongo.ErrNoDocuments {
		return c.newSession(ctx)
	} else if err != nil {
		panic(err)
	}
	c.Session = new(Session)
	c.Session.SessUid = model.MustGetObjectID(session.Id)
	if session.SessionUser != nil {
		c.Session.HasAuthUserUid = true
		c.Session.AuthUserUid = model.MustGetObjectID(session.SessionUser)
		user := new(model.User)
		if err = c.Server.mongodb.Collection("user").FindOne(ctx,
			bson.D{{"_id", session.SessionUser}},
			mongo_options.FindOne().SetProjection(bson.D{{"_id", 1}, {"username", 1}}),
		).Decode(user); err != nil {
			panic(err)
		}
		c.Session.AuthUserUserName = user.GetUsername()
	}
	return nil
}
