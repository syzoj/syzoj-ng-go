package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	model_user "github.com/syzoj/syzoj-ng-go/app/model/user"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type RegisterResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason"`
}

func (srv *ApiServer) HandleAuthRegister(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	req := new(RegisterRequest)
	if err := decoder.Decode(req); err != nil {
		srv.BadRequest(w, err)
		return
	}

	userId, err := util.GenerateUUID()
	if err != nil {
		panic(errors.Wrap(err, "Failed to generate UUID"))
	}
	authInfo := model_user.PasswordAuth(req.Password)
	authInfoJson, err := json.Marshal(authInfo)
	if err != nil {
		panic(err)
	}
	if _, err := srv.db.Exec("INSERT INTO users (id, user_name, auth_info, can_login) VALUES ($1, $2, $3, true)", userId.ToBytes(), req.UserName, authInfoJson); err != nil {
		if sqlErr, ok := err.(*pq.Error); ok {
			if sqlErr.Code == "23505" && sqlErr.Constraint == "users_user_name_unique" {
				srv.SuccessWithError(w, DuplicateUserNameError)
				return
			}
		}
		panic(err)
	}

	srv.Success(w, RegisterResponse{Success: true})
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct{}

func (srv *ApiServer) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	sess := srv.GetSession(r)
	reqDecoder := json.NewDecoder(r.Body)
	var req LoginRequest
	if err := reqDecoder.Decode(&req); err != nil {
		srv.BadRequest(w, err)
		return
	}

	if sess.IsLoggedIn() {
		srv.SuccessWithError(w, AlreadyLoggedInError)
		return
	}
	row := srv.db.QueryRow("SELECT id, auth_info, can_login FROM users WHERE user_name = $1", req.UserName)
	var userIdBytes []byte
	var authInfoJson []byte
	var canLogin bool
	if err := row.Scan(&userIdBytes, &authInfoJson, &canLogin); err != nil {
		if err == sql.ErrNoRows {
			srv.SuccessWithError(w, UnknownUsernameError)
			return
		}
		panic(err)
	}
	userId, err := util.UUIDFromBytes(userIdBytes)
	if err != nil {
		panic(err)
	}
	var authInfo model_user.UserAuthInfo
	err = json.Unmarshal(authInfoJson, &authInfo)
	if err != nil {
		panic(err)
	}
	if !canLogin {
		srv.SuccessWithError(w, CannotLoginError)
		return
	}
	if authInfo.UseTwoFactor {
		srv.SuccessWithError(w, TwoFactorNotSupportedError)
		return
	}
	err = authInfo.PasswordInfo.Verify(req.Password)
	if err != nil {
		srv.SuccessWithError(w, PasswordIncorrectError)
		return
	}

	sess.AuthUserId = userId
	srv.SaveSession(r, w, sess)
	srv.Success(w, LoginResponse{})
}
