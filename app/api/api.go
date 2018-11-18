package api

import (
	"errors"
	"log"
	"encoding/json"
	"net/http"
	"database/sql"
	"github.com/go-redis/redis"
)

type ApiServer struct {
	db *sql.DB
	redis *redis.Client
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

var ApiEndpointNotFoundError = errors.New("API endpoint not found")
// Internal error
var InvalidAuthUserIdError = errors.New("Invalid AuthUserId")

func CreateApiServer(db *sql.DB, redis *redis.Client) (*ApiServer, error) {
	return &ApiServer{
		db: db,
		redis: redis,
	}, nil
}

func (*ApiServer) NotFound(w http.ResponseWriter, e error) {
	response := ErrorResponse{
		Error: e.Error(),
	}
	json, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	http.Error(w, string(json), 404)
}

func (*ApiServer) BadRequest(w http.ResponseWriter, e error) {
	response := ErrorResponse{
		Error: e.Error(),
	}
	json, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	http.Error(w, string(json), 400)
}

func (*ApiServer) InternalServerError(w http.ResponseWriter, e error) {
	log.Println("Error handling http request:", e)
	response := ErrorResponse{
		Error: e.Error(),
	}
	json, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	http.Error(w, string(json), 500)
}

func (*ApiServer) Success(w http.ResponseWriter, d interface{}) {
	response := SuccessResponse{
		Data: d,
	}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		log.Println("Response encoding failed:", err)
	}
}

func (srv *ApiServer) HandleCatchAll(w http.ResponseWriter, r *http.Request) {
	srv.NotFound(w, ApiEndpointNotFoundError)
}