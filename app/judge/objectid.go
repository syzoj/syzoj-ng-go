package judge

import (
	"encoding/base64"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
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
