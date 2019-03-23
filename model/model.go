package model

import (
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
)

func newId() string {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(b[:])
}

type UserRef string

func NewUserRef() UserRef {
	return UserRef(newId())
}

type ProblemRef string

func NewProblemRef() ProblemRef {
	return ProblemRef(newId())
}

type SubmissionRef string

func NewSubmissionRef() SubmissionRef {
	return SubmissionRef(newId())
}

func (m *UserAuth) Value() (driver.Value, error) {
	return proto.Marshal(m)
}

func (m *UserAuth) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	if b, ok := v.([]byte); ok {
		return proto.Unmarshal(b, m)
	}
	return errors.New(fmt.Sprintf("Cannot scan %T into protobuf message", v))
}
