package util

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"errors"

	"github.com/gogo/protobuf/proto"
)

var ErrProtoSqlNil = errors.New("Cannot scan nil into ProtoSql")
var ErrProtoSqlUnknownType = errors.New("Unknown type in ProtoSql")

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

type ProtoSql struct {
	Message proto.Message
}

func (s ProtoSql) Value() (driver.Value, error) {
	return proto.Marshal(s.Message)
}

func (s ProtoSql) Scan(src interface{}) error {
	if src == nil {
		return ErrProtoSqlNil
	}
	if b, ok := src.([]byte); ok {
		return proto.Unmarshal(b, s.Message)
	}
	return ErrProtoSqlUnknownType
}
