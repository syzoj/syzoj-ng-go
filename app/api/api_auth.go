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

func (srv *ApiServer) HandleAuthRegister(w http.ResponseWriter, r *http.Request) {
	var err error
	var sessId uuid.UUID
	var sess *session.Session
	defer func() {
		if err != nil {
			writeError(w, r, err, sess)
		}
	}()
	if sessId, sess, err = srv.ensureSession(w, r); err != nil {
		return
	}

	var req RegisterRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	var userId uuid.UUID
	if userId, err = srv.authService.RegisterUser(req.UserName, req.Password); err != nil {
		return
	}

	sess.AuthUserId = userId
	sess.UserName = req.UserName
	if err = srv.sessService.UpdateSession(sessId, sess); err != nil {
		return
	}
	writeResponse(w, nil, sess)
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (srv *ApiServer) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	var err error
	var sessId uuid.UUID
	var sess *session.Session
	defer func() {
		if err != nil {
			writeError(w, r, err, sess)
		}
	}()
	if sessId, sess, err = srv.ensureSession(w, r); err != nil {
		return
	}
	var req LoginRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	if sess.AuthUserId != defaultUserId {
		err = AlreadyLoggedInError
		return
	}
	var userId uuid.UUID
	if userId, err = srv.authService.LoginUser(req.UserName, req.Password); err != nil {
		return
	}
	sess.AuthUserId = userId
	sess.UserName = req.UserName
	if err = srv.sessService.UpdateSession(sessId, sess); err != nil {
		return
	}

	writeResponse(w, nil, sess)
}

func (srv *ApiServer) HandleAuthLogout(w http.ResponseWriter, r *http.Request) {
	var err error
	var sessId uuid.UUID
	var sess *session.Session
	defer func() {
		if err != nil {
			writeError(w, r, err, sess)
		}
	}()
	if sessId, sess, err = srv.ensureSession(w, r); err != nil {
		return
	}
	if sess.AuthUserId == defaultUserId {
		err = NotLoggedInError
		return
	}
	sess.AuthUserId = defaultUserId
	sess.UserName = ""
	if err = srv.sessService.UpdateSession(sessId, sess); err != nil {
		return
	}

	writeResponse(w, nil, sess)
}
