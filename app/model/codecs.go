package model

import (
	"errors"
	"reflect"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/timestamp"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var objectIDType = reflect.TypeOf(primitive.ObjectID{})
var timeType = reflect.TypeOf(time.Time{})
var durationType = reflect.TypeOf(time.Duration(0))

type objectIDCodec struct{}

func (objectIDCodec) EncodeValue(c bsoncodec.EncodeContext, w bsonrw.ValueWriter, v reflect.Value) error {
	x := v.Interface().(*ObjectID)
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

type anyCodec struct{}

var ErrInvalidAny = errors.New("Invalid any data")

func (anyCodec) EncodeValue(c bsoncodec.EncodeContext, w bsonrw.ValueWriter, v reflect.Value) error {
	x := v.Interface().(*any.Any)
	var dany ptypes.DynamicAny
	if err := ptypes.UnmarshalAny(x, &dany); err != nil {
		return err
	}
	enc, err := c.LookupEncoder(reflect.TypeOf(dany.Message))
	if err != nil {
		return err
	}
	doc, err := w.WriteDocument()
	if err != nil {
		return err
	}
	el, err := doc.WriteDocumentElement("_type")
	if err != nil {
		return err
	}
	if err = el.WriteString(x.TypeUrl); err != nil {
		return err
	}
	el, err = doc.WriteDocumentElement("_val")
	if err != nil {
		return err
	}
	if err = enc.EncodeValue(c, el, reflect.ValueOf(dany.Message)); err != nil {
		return err
	}
	if err = doc.WriteDocumentEnd(); err != nil {
		return err
	}
	return nil
}

func (anyCodec) DecodeValue(c bsoncodec.DecodeContext, r bsonrw.ValueReader, v reflect.Value) error {
	doc, err := r.ReadDocument()
	if err != nil {
		return err
	}
	key, val, err := doc.ReadElement()
	if err != nil {
		return err
	} else if key != "_type" {
		return ErrInvalidAny
	}
	typeUrl, err := val.ReadString()
	if err != nil {
		return err
	}
	x := new(any.Any)
	x.TypeUrl = typeUrl
	msg, err := ptypes.Empty(x)
	if err != nil {
		return err
	}
	dec, err := c.LookupDecoder(reflect.ValueOf(msg).Elem().Type())
	if err != nil {
		return err
	}
	key, val, err = doc.ReadElement()
	if err != nil {
		return err
	} else if key != "_val" {
		return ErrInvalidAny
	}
	if err = dec.DecodeValue(c, val, reflect.ValueOf(msg).Elem()); err != nil {
		return err
	}
	x, err = ptypes.MarshalAny(msg)
	if err != nil {
		return err
	}
	v.Set(reflect.ValueOf(x))
	return nil
}

// Register registers the codecs.
func Register(r *bsoncodec.RegistryBuilder) *bsoncodec.RegistryBuilder {
	return r.RegisterCodec(reflect.TypeOf(&ObjectID{}), objectIDCodec{}).
		RegisterCodec(reflect.TypeOf(&timestamp.Timestamp{}), timestampCodec{}).
		RegisterCodec(reflect.TypeOf(&duration.Duration{}), durationCodec{}).
		RegisterCodec(reflect.TypeOf(&any.Any{}), anyCodec{})
}
