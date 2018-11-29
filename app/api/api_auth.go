package api

import (
	"encoding/json"

	"database/sql"

	"github.com/google/uuid"
	"github.com/lib/pq"
	model_user "github.com/syzoj/syzoj-ng-go/app/model/user"
)

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type RegisterResponse struct{}

func HandleAuthRegister(cxt *ApiContext) ApiResponse {
	var req RegisterRequest
	if err := cxt.ReadBody(&req); err != nil {
		return err
	}
	if err := UseTx(cxt); err != nil {
		return err
	}
	authInfo := model_user.PasswordAuth(req.Password)
	authInfoJson, err := json.Marshal(authInfo)
	if err != nil {
		panic(err)
	}
	userId := uuid.New()
	if _, err := cxt.tx.Exec("INSERT INTO users (id, name, auth_info, can_login) VALUES ($1, $2, $3, true)", userId[:], req.UserName, authInfoJson); err != nil {
		if sqlErr, ok := err.(*pq.Error); ok {
			if sqlErr.Code == "23505" && sqlErr.Constraint == "users_name_unique" {
				return DuplicateUserNameError
			}
		}
		panic(err)
	}

	DoneTx(cxt)
	return Success(nil)
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct{}

func HandleAuthLogin(cxt *ApiContext) ApiResponse {
	var req LoginRequest
	if err := cxt.ReadBody(&req); err != nil {
		return err
	}
	if cxt.sess.IsLoggedIn() {
		return AlreadyLoggedInError
	}
	if err := UseTx(cxt); err != nil {
		return err
	}
	row := cxt.tx.QueryRow("SELECT id, auth_info, can_login FROM users WHERE name = $1", req.UserName)
	var userId uuid.UUID
	var authInfoJson []byte
	var canLogin bool
	if err := row.Scan(&userId, &authInfoJson, &canLogin); err != nil {
		if err == sql.ErrNoRows {
			return UnknownUsernameError
		}
		panic(err)
	}
	var authInfo model_user.UserAuthInfo
	err := json.Unmarshal(authInfoJson, &authInfo)
	if err != nil {
		panic(err)
	}
	if !canLogin {
		return CannotLoginError
	}
	if authInfo.UseTwoFactor {
		return TwoFactorNotSupportedError
	}
	err = authInfo.PasswordInfo.Verify(req.Password)
	if err != nil {
		return PasswordIncorrectError
	}

	cxt.sess.AuthUserId = userId
	cxt.sess.Save()
	DoneTx(cxt)
	return Success(nil)
}
