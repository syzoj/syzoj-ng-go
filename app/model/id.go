package model

import (
	"crypto/rand"
	"encoding/base64"
)

func NewID() string {
	var v [15]byte
	rand.Read(v[:])
	return base64.URLEncoding.EncodeToString(v[:])
}
