package model

import (
	"encoding/base64"
	"errors"

	"github.com/golang/protobuf/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidObjectID = errors.New("Invalid ObjectID")

func EncodeObjectID(id primitive.ObjectID) (string, error) {
	return base64.URLEncoding.EncodeToString(id[:]), nil
}

func DecodeObjectID(id string) (primitive.ObjectID, error) {
	var v primitive.ObjectID
	n, err := base64.URLEncoding.Decode(v[:], []byte(id))
	if err != nil || n != 12 {
		return primitive.ObjectID{}, ErrInvalidObjectID
	}
	return v, nil
}

func MustDecodeObjectID(id string) primitive.ObjectID {
	v, err := DecodeObjectID(id)
	if err != nil {
		panic(err)
	}
	return v
}

func ObjectIDProto(id primitive.ObjectID) *ObjectID {
	s, _ := EncodeObjectID(id)
	return &ObjectID{Id: proto.String(s)}
}

func NewObjectIDProto() *ObjectID {
	return ObjectIDProto(primitive.NewObjectID())
}

func GetObjectID(o *ObjectID) (primitive.ObjectID, error) {
	if o == nil {
		return primitive.ObjectID{}, ErrInvalidObjectID
	}
	return DecodeObjectID(*o.Id)
}

func MustGetObjectID(o *ObjectID) primitive.ObjectID {
	if o == nil {
		panic(ErrInvalidObjectID)
	}
	v, err := DecodeObjectID(*o.Id)
	if err != nil {
		panic(err)
	}
	return v
}
