package api

import (
	"encoding/json"

	model_user "github.com/syzoj/syzoj-ng-go/app/model/user"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type UserInfoResponse struct {
	UserId    util.UUID `json:"user_id"`
	Biography string    `json:"biography"`
}

func HandleUserInfo(cxt *ApiContext) *ApiError {
	if !cxt.sess.IsLoggedIn() {
		return NotLoggedInError
	}
	row := cxt.s.db.QueryRow("SELECT profile_info FROM users WHERE Id=$1", cxt.sess.AuthUserId.ToBytes())
	var userInfoData []byte
	if err := row.Scan(&userInfoData); err != nil {
		panic(err)
	}
	var userInfo model_user.UserProfileInfo
	if err := json.Unmarshal(userInfoData, &userInfo); err != nil {
		panic(err)
	}

	cxt.resp = UserInfoResponse{UserId: cxt.sess.AuthUserId, Biography: userInfo.Biography}
	return nil
}
