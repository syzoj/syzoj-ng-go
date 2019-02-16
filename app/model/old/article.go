package model

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Article struct {
	Id           primitive.ObjectID `bson:"_id"`
	Title        string             `bson:"title"`
	Owner        primitive.ObjectID `bson:"owner"`
	Text         string             `bson:"text"`
	Reply        []Reply            `bson:"reply"`
	CreateTime   time.Time          `bson:"create_time"`
	LastEditTime time.Time          `bson:"last_edit_time"`
}

type Reply struct {
	Owner        primitive.ObjectID `bson:"owner"`
	CreateTime   time.Time          `bson:"create_time"`
	LastEditTime time.Time          `bson:"last_edit_time"`
	Text         string             `bson:"text"`
}
