package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/syzoj/syzoj-ng-go/app/session"
)

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (srv *ApiServer) HandleAuthRegister(w http.ResponseWriter, r *http.Request, sessId uuid.UUID, sess *session.Session) ApiError {
	var err error
	var req RegisterRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError(err)
	}
	var userId uuid.UUID
	if userId, err = srv.authService.RegisterUser(req.UserName, req.Password); err != nil {
		return userError(err)
	}
	sess.AuthUserId = userId
	sess.UserName = req.UserName
	if err = srv.updateSession(sessId, sess); err != nil {
		return err.(ApiError)
	}
	writeResponse(w, nil, sess)
	return nil
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (srv *ApiServer) HandleAuthLogin(w http.ResponseWriter, r *http.Request, sessId uuid.UUID, sess *session.Session) ApiError {
	var err error
	var req LoginRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return badRequestError(err)
	}
	if err = requireLogin(sess); err != nil {
		return err.(ApiError)
	}
	var userId uuid.UUID
	if userId, err = srv.authService.LoginUser(req.UserName, req.Password); err != nil {
		return userError(err)
	}
	sess.AuthUserId = userId
	sess.UserName = req.UserName
	if err = srv.updateSession(sessId, sess); err != nil {
		return err.(ApiError)
	}
	writeResponse(w, nil, sess)
	return nil
}

func (srv *ApiServer) HandleAuthLogout(w http.ResponseWriter, r *http.Request, sessId uuid.UUID, sess *session.Session) ApiError {
	var err ApiError
	if err = requireLogin(sess); err != nil {
		return err.(ApiError)
	}
	sess.AuthUserId = defaultUserId
	sess.UserName = ""
	if err = srv.updateSession(sessId, sess); err != nil {
		return err.(ApiError)
	}
	writeResponse(w, nil, sess)
	return nil
}
