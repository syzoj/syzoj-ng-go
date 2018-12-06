package auth

import (
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserAuthInfo struct {
	UseTwoFactor bool             `json:"use_two_factor"`
	PasswordInfo UserPasswordInfo `json:"password_info"`
}

type UserPasswordInfo struct {
	Type int         `json:"type"`
	Data interface{} `json:"data"`
}

var InvalidAuthInfoError = errors.New("Invalid auth info")

type UserProfileInfo struct {
	Biography string `json:"biography"`
}

func BcryptPassword(password string) UserPasswordInfo {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		panic(err)
	}
	hash_base64 := base64.StdEncoding.EncodeToString(hash)
	return UserPasswordInfo{
		Type: 1,
		Data: hash_base64,
	}
}

func PasswordAuth(password string) UserAuthInfo {
	info := BcryptPassword(password)
	return UserAuthInfo{
		UseTwoFactor: false,
		PasswordInfo: info,
	}
}

func (info UserPasswordInfo) Verify(password string) error {
	switch info.Type {
	case 1:
		hash_base64 := info.Data.(string)
		hash, err := base64.StdEncoding.DecodeString(hash_base64)
		if err != nil {
			return InvalidAuthInfoError
		}

		err = bcrypt.CompareHashAndPassword(hash, []byte(password))
		return err
	}
	return InvalidAuthInfoError
}
