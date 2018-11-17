package api

import (
	"github.com/syzoj/syzoj-ng-go/app/util"
	"net/http"
)

type UserInfoResponse struct {
	LoggedIn bool `json:"logged_in"`
	UserId util.UUID `json:"user_id"`
}
func (srv *ApiServer) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	session := srv.GetSession(r)
	if !session.LoggedIn {
		srv.Success(w, UserInfoResponse{LoggedIn: false})
		return
	}

	userId := session.AuthUserId
	srv.Success(w, UserInfoResponse{LoggedIn: true, UserId: userId,})
}