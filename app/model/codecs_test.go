package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

type Message struct {
	Oid       *ObjectID
	T         *timestamp.Timestamp
	D         *duration.Duration
	CamelCase string `json:"camel_case" bson:"camel_case"`
}

func Test_X(t *testing.T) {
	builder := bsoncodec.NewRegistryBuilder()
	bsoncodec.DefaultValueEncoders{}.RegisterDefaultEncoders(builder)
	bsoncodec.DefaultValueDecoders{}.RegisterDefaultDecoders(builder)
	Register(builder)
	registry := builder.Build()

	m := Message{Oid: &ObjectID{Id: proto.String("XXXXXXXXXXXXXXXX")}, D: ptypes.DurationProto(time.Second * 3), CamelCase: "X"}
	m.T, _ = ptypes.TimestampProto(time.Now())
	b, err := bson.MarshalWithRegistry(registry, m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Original message: %+v\n", m)

	var i interface{}
	bson.Unmarshal(b, &i)
	j, err := bson.MarshalExtJSON(i, false, false)
	if err != nil {
		panic(err)
	}
	fmt.Printf("ExtJSON: %s\n", string(j))

	m = Message{}
	err = bson.UnmarshalWithRegistry(registry, b, &m)
	if err != nil {
		panic(err)
	}
	fmt.Printf("BSON message: %+v\n", b)
	fmt.Printf("Unmarshaled message: %+v\n", m)
}
