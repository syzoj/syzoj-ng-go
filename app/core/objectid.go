package core

import (
	"encoding/base64"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func EncodeObjectID(id primitive.ObjectID) string {
	return base64.URLEncoding.EncodeToString(id[:])
}

func DecodeObjectID(id string) (res primitive.ObjectID) {
	n, err := base64.URLEncoding.Decode(res[:], []byte(id))
	if err != nil || n != 12 {
		panic("Invalid ObjectID string")
	}
	return
}

func DecodeObjectIDOK(id string) (res primitive.ObjectID, ok bool) {
	n, err := base64.URLEncoding.Decode(res[:], []byte(id))
	if err != nil || n != 12 {
		ok = false
		return
	}
	ok = true
	return
}
