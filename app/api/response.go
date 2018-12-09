package api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log = logrus.StandardLogger()

type ApiErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
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
	Data interface{} `json:"data"`
}

func writeResponse(w http.ResponseWriter, data interface{}) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(ApiSuccessResponse{data}); err != nil {
		log.WithField("error", err).Warning("Failed to write response")
	}
}
