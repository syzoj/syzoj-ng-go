package model

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	UserName     string             `bson:"username,omitempty"`
	RegisterTime time.Time          `bson:"register_time,omitempty"`
	Auth         UserAuth           `bson:"auth,omitempty"`
}

type UserAuth struct {
	Method   int64  `bson:"method",omitempty` // 1: password login
	Password []byte `bson:"password",omitempty`
}
