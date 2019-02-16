package model

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Session struct {
	Id           primitive.ObjectID `bson:"_id"`
	SessionToken string             `bson:"session_token,omitempty"`
	SessionUser  primitive.ObjectID `bson:"session_user,omitempty"`
}
