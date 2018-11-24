package util

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
)

type UUID [16]byte

func GenerateUUID() (uuid UUID, err error) {
	_, err = rand.Read(uuid[:])
	return
}

func (uuid UUID) String() string {
	var buf [36]byte

	hex.Encode(buf[0:8], uuid[0:4])
	hex.Encode(buf[9:13], uuid[4:6])
	hex.Encode(buf[14:18], uuid[6:8])
	hex.Encode(buf[19:23], uuid[8:10])
	hex.Encode(buf[24:36], uuid[10:16])
	buf[8] = '-'
	buf[13] = '-'
	buf[18] = '-'
	buf[23] = '-'

	return string(buf[0:36])
}

func (uuid UUID) ToBytes() []byte {
	return uuid[:]
}

func (uuid UUID) MarshalBinary() ([]byte, error) {
	return uuid[:], nil
}

func (uuid *UUID) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return errors.New("uuid: Byte array length is not 16")
	}
	copy(uuid[:], data)
	return nil
}

func (uuid UUID) MarshalText() ([]byte, error) {
	return []byte(uuid.String()), nil
}

func (uuid *UUID) UnmarshalText(text []byte) (err error) {
	*uuid, err = ParseUUID(text)
	return
}

func UUIDFromBytes(data []byte) (uuid UUID, err error) {
	if len(data) != 16 {
		err = errors.New("uuid: Byte array rray length is not 16")
		return
	}
	copy(uuid[:], data[0:16])
	return
}

func ParseUUIDString(str string) (UUID, error) {
	return ParseUUID([]byte(str))
}

func ParseUUID(buf []byte) (uuid UUID, err error) {
	if len(buf) != 36 || buf[8] != '-' || buf[13] != '-' || buf[18] != '-' || buf[23] != '-' {
		return UUID{}, errors.New("Invalid UUID")
	}

	if _, err := hex.Decode(uuid[0:4], buf[0:8]); err != nil {
		return UUID{}, errors.New("Invalid UUID")
	}

	if _, err := hex.Decode(uuid[4:6], buf[9:13]); err != nil {
		return UUID{}, errors.New("Invalid UUID")
	}

	if _, err := hex.Decode(uuid[6:8], buf[14:18]); err != nil {
		return UUID{}, errors.New("Invalid UUID")
	}

	if _, err := hex.Decode(uuid[8:10], buf[19:23]); err != nil {
		return UUID{}, errors.New("Invalid UUID")
	}

	if _, err := hex.Decode(uuid[10:16], buf[24:36]); err != nil {
		return UUID{}, errors.New("Invalid UUID")
	}

	return
}

func (uuid *UUID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	var err error
	if *uuid, err = ParseUUIDString(s); err != nil {
		return err
	}

	return nil
}

func (uuid UUID) MarshalJSON() ([]byte, error) {
	s := uuid.String()
	return json.Marshal(s)
}
