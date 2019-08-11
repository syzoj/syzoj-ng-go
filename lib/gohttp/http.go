package http

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

func SendJSON(w http.ResponseWriter, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.WithError(err).Error("SendJSON: Failed to marshal data")
		return
	}
}

func SendError(w http.ResponseWriter, err string) {
	SendJSON(w, struct {
		Error string `json:"error"`
	}{Error: err})
}

func SendInternalError(w http.ResponseWriter, err error) {
	log.WithError(err).Error("Internal server error")
	http.Error(w, "Internal server error", 500)
}

func SendConflict(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 409)
}

func ReadBodyJSON(r *http.Request, val interface{}) error {
	mediatype, _, err := mime.ParseMediaType(string(r.Header.Get("Content-Type")))
	if err != nil {
		return err
	}
	if mediatype != "application/json" {
		return fmt.Errorf("Invalid media type: %s", mediatype)
	}
	return json.NewDecoder(r.Body).Decode(val)
}

func BadRequest(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 400)
}

func NotFound(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), 404)
}
