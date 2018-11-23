package api

import (
	"encoding/json"
	"net/http"

	model_user "github.com/syzoj/syzoj-ng-go/app/model/user"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type UserInfoResponse struct {
	UserId    util.UUID `json:"user_id"`
	Biography string    `json:"biography"`
}

func (srv *ApiServer) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	sess := srv.GetSession(r)
	if !sess.IsLoggedIn() {
		srv.SuccessWithError(w, NotLoggedInError)
		return
	}

	userId := sess.AuthUserId
	row := srv.db.QueryRow("SELECT user_profile_info FROM users WHERE Id=$1", userId.ToBytes())
	var userInfoData []byte
	if err := row.Scan(&userInfoData); err != nil {
		panic(err)
	}
	var userInfo model_user.UserProfileInfo
	if err := json.Unmarshal(userInfoData, &userInfo); err != nil {
		panic(err)
	}

	srv.Success(w, UserInfoResponse{UserId: userId, Biography: userInfo.Biography})
}
