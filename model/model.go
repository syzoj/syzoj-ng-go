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

type DeviceRef string

func NewDeviceRef() DeviceRef {
	return DeviceRef(newId())
}

type ProblemRef string

func NewProblemRef() ProblemRef {
	return ProblemRef(newId())
}

type ProblemSourceRef string

func NewProblemSourceRef() ProblemSourceRef {
	return ProblemSourceRef(newId())
}

type ProblemJudgerRef string

func NewProblemJudgerRef() ProblemJudgerRef {
	return ProblemJudgerRef(newId())
}

type ProblemStatementRef string

func NewProblemStatementRef() ProblemStatementRef {
	return ProblemStatementRef(newId())
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

func (m *DeviceInfo) Value() (driver.Value, error) {
	return proto.Marshal(m)
}

func (m *DeviceInfo) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	if b, ok := v.([]byte); ok {
		return proto.Unmarshal(b, m)
	}
	return errors.New(fmt.Sprintf("Cannot scan %T into protobuf message", v))
}
