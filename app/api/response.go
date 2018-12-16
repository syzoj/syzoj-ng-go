package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/app/session"
)

var log = logrus.StandardLogger()

type ApiErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	log.Infof("Error handling request %s: %s", r.URL, err)
	switch v := err.(type) {
	case *ApiError:
		w.WriteHeader(v.Code)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(ApiErrorResponse{v.Message}); err != nil {
			log.WithField("error", err).Warning("Failed to write error")
		}
	default:
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(ApiErrorResponse{v.Error()}); err != nil {
			log.WithField("error", err).Warning("Failed to write error")
		}
	}
}

type ApiSuccessResponse struct {
	Data    interface{}     `json:"data"`
	Session SessionResponse `json:"session"`
}
type SessionResponse struct {
	UserId   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
}

func writeResponse(w http.ResponseWriter, data interface{}) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(ApiSuccessResponse{Data: data}); err != nil {
		log.WithField("error", err).Warning("Failed to write response")
	}
}

func writeResponseWithSession(w http.ResponseWriter, data interface{}, sess *session.Session) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(ApiSuccessResponse{Data: data, Session: SessionResponse{
		UserId:   sess.AuthUserId,
		UserName: sess.UserName,
	}}); err != nil {
		log.WithField("error", err).Warning("Failed to write response")
	}
}
