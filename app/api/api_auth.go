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
	defer func() {
		if err != nil {
			writeError(w, err)
		}
	}()
	var sessId uuid.UUID
	var sess *session.Session
	if sessId, sess, err = srv.ensureSession(r); err != nil {
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
	if err = srv.sessService.UpdateSession(sessId, sess); err != nil {
		return
	}
	writeResponse(w, nil)
}

type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func (srv *ApiServer) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			writeError(w, err)
		}
	}()
	var sessId uuid.UUID
	var sess *session.Session
	if sessId, sess, err = srv.ensureSession(r); err != nil {
		return
	}

	var req LoginRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	var userId uuid.UUID
	if userId, err = srv.authService.LoginUser(req.UserName, req.Password); err != nil {
		return
	}
	sess.AuthUserId = userId
	if err = srv.sessService.UpdateSession(sessId, sess); err != nil {
		return
	}

	writeResponse(w, nil)
}
