package model

import (
    "database/sql/driver"
    "fmt"
    "errors"

    "github.com/golang/protobuf/proto"
)

type UserRef string

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
