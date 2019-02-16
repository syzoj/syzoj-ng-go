package model

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mongodb/mongo-go-driver/bson/bsoncodec"
	"github.com/mongodb/mongo-go-driver/bson/bsonrw"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

var objectIDType = reflect.TypeOf(primitive.ObjectID{})
var timeType = reflect.TypeOf(time.Time{})
var durationType = reflect.TypeOf(time.Duration(0))

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

// MarshalJSONPB implements the jsonpb.JSONPBMarshaler interface.
func (o *ObjectID) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	if o.Id == nil {
		return nil, ErrInvalidObjectID
	}
	_, err := DecodeObjectID(*o.Id)
	if err != nil {
		return nil, ErrInvalidObjectID
	}
	return []byte(`"` + *o.Id + `"`), nil
}

// UnmarshalJSONPB implements the jsonpb.JSONPBUnmarshaler interface.
func (o *ObjectID) UnmarshalJSONPB(m *jsonpb.Unmarshaler, b []byte) error {
	var (
		s   string
		err error
	)
	if err = json.Unmarshal(b, &s); err != nil {
		return err
	}
	if _, err = DecodeObjectID(s); err != nil {
		return err
	}
	o.Id = proto.String(s)
	return nil
}

type objectIDCodec struct{}

func (objectIDCodec) EncodeValue(c bsoncodec.EncodeContext, w bsonrw.ValueWriter, v reflect.Value) error {
	x := v.Interface().(*ObjectID)
	if x == nil {
		return ErrInvalidObjectID
	}
	objectID, err := DecodeObjectID(*x.Id)
	if err != nil {
		return err
	}
	enc, err := c.LookupEncoder(objectIDType)
	if err != nil {
		return err
	}
	return enc.EncodeValue(c, w, reflect.ValueOf(objectID))
}
func (objectIDCodec) DecodeValue(c bsoncodec.DecodeContext, r bsonrw.ValueReader, v reflect.Value) error {
	dec, err := c.LookupDecoder(objectIDType)
	if err != nil {
		return err
	}
	var objectID primitive.ObjectID
	if err = dec.DecodeValue(c, r, reflect.ValueOf(&objectID).Elem()); err != nil {
		return err
	}
	v.Set(reflect.ValueOf(ObjectIDProto(objectID)))
	return nil
}

type timestampCodec struct{}

func (timestampCodec) EncodeValue(c bsoncodec.EncodeContext, w bsonrw.ValueWriter, v reflect.Value) error {
	x := v.Interface().(*timestamp.Timestamp)
	t, err := ptypes.Timestamp(x)
	if err != nil {
		return err
	}
	enc, err := c.LookupEncoder(timeType)
	if err != nil {
		return err
	}
	return enc.EncodeValue(c, w, reflect.ValueOf(t))
}
func (timestampCodec) DecodeValue(c bsoncodec.DecodeContext, r bsonrw.ValueReader, v reflect.Value) error {
	dec, err := c.LookupDecoder(timeType)
	if err != nil {
		return err
	}
	var t time.Time
	if err = dec.DecodeValue(c, r, reflect.ValueOf(&t).Elem()); err != nil {
		return err
	}
	tproto, err := ptypes.TimestampProto(t.In(time.UTC))
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(tproto))
	return nil
}

type durationCodec struct{}

func (durationCodec) EncodeValue(c bsoncodec.EncodeContext, w bsonrw.ValueWriter, v reflect.Value) error {
	x := v.Interface().(*duration.Duration)
	t, err := ptypes.Duration(x)
	if err != nil {
		return err
	}
	enc, err := c.LookupEncoder(durationType)
	if err != nil {
		return err
	}
	return enc.EncodeValue(c, w, reflect.ValueOf(t))
}
func (durationCodec) DecodeValue(c bsoncodec.DecodeContext, r bsonrw.ValueReader, v reflect.Value) error {
	dec, err := c.LookupDecoder(durationType)
	if err != nil {
		return err
	}
	var t time.Duration
	if err = dec.DecodeValue(c, r, reflect.ValueOf(&t).Elem()); err != nil {
		return err
	}
	tproto := ptypes.DurationProto(t)
	v.Set(reflect.ValueOf(tproto))
	return nil
}

// Register registers the codecs.
func Register(r *bsoncodec.RegistryBuilder) *bsoncodec.RegistryBuilder {
	return r.RegisterCodec(reflect.TypeOf(&ObjectID{}), objectIDCodec{}).
		RegisterCodec(reflect.TypeOf(&timestamp.Timestamp{}), timestampCodec{}).
		RegisterCodec(reflect.TypeOf(&duration.Duration{}), durationCodec{})
}
