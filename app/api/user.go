package api

import (
	"errors"
	"encoding/base64"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthInfo struct {
	UseTwoFactor bool
	PasswordInfo UserPasswordInfo
}

type UserPasswordInfo struct {
	Type int `json:"type"`
	Data interface{} `json:"data"`
}

func BcryptPassword(password string) (UserPasswordInfo, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return UserPasswordInfo{}, err
	}
	hash_base64 := base64.StdEncoding.EncodeToString(hash)
	return UserPasswordInfo{
		Type: 1,
		Data: hash_base64,
	}, nil
}

func PasswordAuth(password string) (UserAuthInfo, error) {
	info, err := BcryptPassword(password)
	if err != nil {
		return UserAuthInfo{}, err
	}
	return UserAuthInfo {
		UseTwoFactor: false,
		PasswordInfo: info,
	}, nil
}

func (info UserPasswordInfo) Verify(password string) (bool, error) {
	switch info.Type {
	case 1:
		hash_base64, ok := info.Data.(string)
		if !ok {
			return false, errors.New("Invalid password info")
		}

		hash, err := base64.StdEncoding.DecodeString(hash_base64)
		if err != nil {
			return false, err
		}

		err = bcrypt.CompareHashAndPassword(hash, []byte(password))
		return err == nil, nil
	}
	
	return false, errors.New("Invalid password info")
}