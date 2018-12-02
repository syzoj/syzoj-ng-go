package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type ApiErrorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	switch v := err.(type) {
	case *ApiError:
		w.WriteHeader(v.Code)
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(ApiErrorResponse{v.Message}); err != nil {
			log.Println("Warning: failed to write error: ", err)
		}
	default:
		encoder := json.NewEncoder(w)
		if err := encoder.Encode(ApiErrorResponse{v.Error()}); err != nil {
			log.Println("Warning: failed to write error: ", err)
		}
	}
}

type ApiSuccessResponse struct {
	Data interface{} `json:"data"`
}

func writeResponse(w http.ResponseWriter, data interface{}) {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(ApiSuccessResponse{data}); err != nil {
		log.Println("Warning: failed to write response: ", err)
	}
}
