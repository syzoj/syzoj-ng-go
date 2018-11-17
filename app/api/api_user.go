package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/syzoj/syzoj-ng-go/app/util"
	"github.com/syzoj/syzoj-ng-go/app/model"
)

type UserInfoResponse struct {
	LoggedIn bool `json:"logged_in"`
	UserId util.UUID `json:"user_id"`
	Biography string `json:"biography"`
}
func (srv *ApiServer) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	session := srv.GetSession(r)
	if !session.LoggedIn {
		srv.Success(w, UserInfoResponse{LoggedIn: false})
		return
	}

	userId := session.AuthUserId
	rows, err := srv.db.Query("SELECT user_profile_info FROM users WHERE Id=$1", userId.ToBytes())
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}
	if !rows.Next() {
		srv.InternalServerError(w, errors.New("Invalid AuthUserId"))
		return
	}
	var userInfoData []byte
	rows.Scan(&userInfoData)
	var userInfo model.UserProfileInfo
	if err := json.Unmarshal(userInfoData, &userInfo); err != nil {
		srv.InternalServerError(w, err)
		return
	}
	
	srv.Success(w, UserInfoResponse{LoggedIn: true, UserId: userId, Biography: userInfo.Biography})
}