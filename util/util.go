package util

import (
	"crypto/rand"
	"encoding/base64"
)

func RandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b)
}

func NewId() string {
	return RandomString(12)
}
