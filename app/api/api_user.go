package api

import (
	"encoding/json"

	"github.com/google/uuid"

	model_user "github.com/syzoj/syzoj-ng-go/app/model/user"
)

type UserInfoResponse struct {
	UserId    uuid.UUID `json:"user_id"`
	Biography string    `json:"biography"`
}

func HandleUserInfo(cxt *ApiContext) ApiResponse {
	if !cxt.sess.IsLoggedIn() {
		return NotLoggedInError
	}
	row := cxt.s.db.QueryRow("SELECT profile_info FROM users WHERE Id=$1", cxt.sess.AuthUserId[:])
	var userInfoData []byte
	if err := row.Scan(&userInfoData); err != nil {
		panic(err)
	}
	var userInfo model_user.UserProfileInfo
	if err := json.Unmarshal(userInfoData, &userInfo); err != nil {
		panic(err)
	}

	return Success(UserInfoResponse{UserId: cxt.sess.AuthUserId, Biography: userInfo.Biography})
}
