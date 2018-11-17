package api

import (
	"github.com/lib/pq"
	"encoding/json"
	"net/http"

	"github.com/syzoj/syzoj-ng-go/app/util"
)

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type RegisterResponse struct {
	Success bool `json:"success"`
	Reason string `json:"reason"`
	UserId util.UUID `json:"user_id"`
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
		srv.InternalServerError(w, err)
		return
	}

	authInfo, err := PasswordAuth(req.Password)
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}

	authInfoJson, err := json.Marshal(authInfo)
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}

	rows, err := srv.db.Query("INSERT INTO users (id, user_name, auth_info, can_login) VALUES ($1, $2, $3, true)", userId.ToBytes(), req.UserName, authInfoJson)
	if err != nil {
		if sqlErr, ok := err.(*pq.Error); ok {
			if sqlErr.Code == "23505" && sqlErr.Constraint == "users_user_name" {
				srv.Success(w, RegisterResponse{Success: false, Reason: "Duplicate user name"})
				return
			}
		}
		srv.InternalServerError(w, err)
		return
	}
	defer rows.Close()
	
	srv.Success(w, RegisterResponse{Success: true, UserId: userId,})
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Success bool `json:"success"`
	Reason string `json:"reason"`
}
func (srv *ApiServer) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	session := srv.GetSession(r)
	decoder := json.NewDecoder(r.Body)
	req := new(LoginRequest)
	if err := decoder.Decode(req); err != nil {
		srv.BadRequest(w, err)
		return
	}

	if session.LoggedIn {
		srv.Success(w, LoginResponse{Success: false, Reason: "Already logged in"})
		return
	}

	rows, err := srv.db.Query("SELECT id, auth_info, can_login FROM users WHERE user_name = $1", req.UserName)
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		srv.Success(w, LoginResponse{Success: false, Reason: "Unknown username"})
		return
	}

	var userIdBytes []byte
	var authInfoJson []byte
	var canLogin bool
	err = rows.Scan(&userIdBytes, &authInfoJson, &canLogin)
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}
	userId, err := util.UUIDFromBytes(userIdBytes)
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}
	var authInfo = UserAuthInfo{
		UseTwoFactor: false,
	}
	err = json.Unmarshal(authInfoJson, &authInfo)
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}

	if !canLogin {
		srv.Success(w, LoginResponse{Success: false, Reason: "Cannot login yet"})
		return
	}

	if authInfo.UseTwoFactor {
		srv.Success(w, LoginResponse{Success: false, Reason: "Two factor auth not supported"})
		return
	}

	ok, err := authInfo.PasswordInfo.Verify(req.Password)
	if err != nil {
		srv.InternalServerError(w, err)
		return
	}
	if !ok {
		srv.Success(w, LoginResponse{Success: false, Reason: "Password incorrect"})
		return
	}

	session.LoggedIn = true
	session.AuthUserId = userId
	if err := srv.SaveSession(r, w, session); err != nil {
		srv.InternalServerError(w, err)
		return
	}
	srv.Success(w, LoginResponse{Success: true})
}