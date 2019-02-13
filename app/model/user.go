package model

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	UserName     string             `bson:"username,omitempty"`
	Email        string             `bson:"email,omitempty"`
	RegisterTime time.Time          `bson:"register_time,omitempty"`
	Auth         UserAuth           `bson:"auth,omitempty"`
}

type UserAuth struct {
	Method   int64  `bson:"method",omitempty` // 1: password login 2: Legacy SYZOJ hash
	Password []byte `bson:"password",omitempty`
}
