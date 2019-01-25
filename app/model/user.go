package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	Xid          *uuid.UUID         `bson:"xid"`
	UserName     *string            `bson:"username,omitempty"`
	RegisterTime time.Time          `bson:"register_time,omitempty"`
	Auth         *UserAuth          `bson:"auth,omitempty"`
}

type UserAuth struct {
	Password []byte `bson:"password",omitempty`
}
