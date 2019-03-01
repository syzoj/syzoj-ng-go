package tool_import

import (
	"context"
	"database/sql"
	"encoding/hex"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/syzoj/syzoj-ng-go/app/model"
)

type user struct {
	Id           int64
	UserName     string
	Password     string
	Email        string
	RegisterTime sql.NullInt64
}

func (i *importer) readUsers(users chan *user) {
	var err error
	var rows *sql.Rows
	if rows, err = i.db.Query("SELECT id, username, password, email, register_time FROM user"); err != nil {
		log.Fatal("Error importing users from MySQL: ", err.Error())
	}
	for rows.Next() {
		u := new(user)
		err = rows.Scan(&u.Id, &u.UserName, &u.Password, &u.Email, &u.RegisterTime)
		if err != nil {
			log.Error("Error reading user: ", err)
			err = nil
		}
		users <- u
	}
	close(users)
}

func (i *importer) writeUsers(users chan *user) {
	var err error
	i.userId = make(map[int64]primitive.ObjectID)
	for user := range users {
		var passmd5 []byte
		passmd5, err = hex.DecodeString(user.Password)
		if err != nil {
			log.WithField("username", user.UserName).Error("Error parsing password")
			err = nil
			continue
		}
		userModel := new(model.User)
		userModel.Id = model.NewObjectIDProto()
		userModel.Username = proto.String(user.UserName)
		userModel.Email = proto.String(user.Email)
		userModel.Auth = &model.UserAuth{
			Method:   proto.Int64(2),
			Password: passmd5,
		}
		if user.RegisterTime.Valid {
			userModel.RegisterTime, _ = ptypes.TimestampProto(time.Unix(user.RegisterTime.Int64, 0))
		}
		if _, err = i.mongodb.Collection("user").InsertOne(context.Background(), userModel); err != nil {
			log.Error("Error inserting user: ", err)
			err = nil
		}
		i.userId[user.Id] = model.MustGetObjectID(userModel.Id)
	}
}
