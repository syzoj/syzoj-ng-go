package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/go-redis/redis"
	model_group "github.com/syzoj/syzoj-ng-go/app/model/group"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type ApiServer struct {
	db    *sql.DB
	redis *redis.Client
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

var ApiEndpointNotFoundError = errors.New("API endpoint not found")
var GroupNotFoundError = errors.New("Group not found")
var CreateProblemsetDeniedError = errors.New("Cannot create problemset")
var NotLoggedInError = errors.New("Not logged in")
var DuplicateGroupNameError = errors.New("Duplicate group name")
var DuplicateUserNameError = errors.New("Duplicate user name")
var DuplicateProblemsetNameError = errors.New("Duplicate problemset name")
var InvalidProblemsetTypeError = errors.New("Invalid or unsupported problemset type")

var AlreadyLoggedInError = errors.New("Already logged in")
var UnknownUsernameError = errors.New("Unknown username")
var CannotLoginError = errors.New("Cannot login yet")
var TwoFactorNotSupportedError = errors.New("Two factor auth not supported")
var PasswordIncorrectError = errors.New("Password incorrect")

// Internal error
var InvalidAuthUserIdError = errors.New("Invalid AuthUserId")

func CreateApiServer(db *sql.DB, redis *redis.Client) (*ApiServer, error) {
	return &ApiServer{
		db:    db,
		redis: redis,
	}, nil
}

func respondWithError(w http.ResponseWriter, e error, c int) {
	response := ErrorResponse{
		Error: e.Error(),
	}
	json, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	http.Error(w, string(json), c)
}

func (*ApiServer) NotFound(w http.ResponseWriter, e error) {
	respondWithError(w, e, 404)
}

func (*ApiServer) Forbidden(w http.ResponseWriter, e error) {
	respondWithError(w, e, 403)
}

func (*ApiServer) BadRequest(w http.ResponseWriter, e error) {
	respondWithError(w, e, 400)
}

func (*ApiServer) SuccessWithError(w http.ResponseWriter, e error) {
	respondWithError(w, e, 200)
}

func (*ApiServer) Success(w http.ResponseWriter, d interface{}) {
	response := SuccessResponse{Data: d}
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		panic("Response encoding failed:" + err.Error())
	}
}

func (srv *ApiServer) HandleCatchAll(w http.ResponseWriter, r *http.Request) {
	srv.NotFound(w, ApiEndpointNotFoundError)
}

func marshalJson(data interface{}) []byte {
	result, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return result
}

func (srv *ApiServer) InternalServerErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				msg := fmt.Sprintf("%#v", err)
				log.Println("Error handling http request:", msg)
				debug.PrintStack()
				response := ErrorResponse{
					Error: msg,
				}
				json, _ := json.Marshal(response)
				http.Error(w, string(json), 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (srv *ApiServer) GetGroupPolicyByName(groupName string) (groupId util.UUID, groupPolicy model_group.GroupPolicy) {
	row := srv.db.QueryRow("SELECT id, policy_info FROM groups WHERE group_name=$1", groupName)
	var idBytes []byte
	var infoBytes []byte
	if err := row.Scan(&idBytes, &infoBytes); err != nil {
		if err == sql.ErrNoRows {
			return
		}
		panic(err)
	}
	groupId, err := util.UUIDFromBytes(idBytes)
	if err != nil {
		panic(err)
	}
	groupProvider := model_group.GetGroupType()
	groupPolicy = groupProvider.GetDefaultGroupPolicy()
	if err := json.Unmarshal(infoBytes, &groupPolicy); err != nil {
		panic(err)
	}
	return
}

func (srv *ApiServer) GetGroupUserRole(groupId util.UUID, groupPolicy model_group.GroupPolicy, userId util.UUID) (userRole model_group.GroupUserRole) {
    if userId == (util.UUID{}) {
        return groupPolicy.GetGuestRole()
    }
	row := srv.db.QueryRow("SELECT role_info FROM group_users WHERE group_id=$1 AND user_id=$2", groupId.ToBytes(), userId.ToBytes())
	var infoBytes []byte
	if err := row.Scan(&infoBytes); err != nil {
		if err == sql.ErrNoRows {
            userRole = groupPolicy.GetRegisteredUserRole()
			return
		}
		panic(err)
	}
	userRole = groupPolicy.GetDefaultRole()
	if err := json.Unmarshal(infoBytes, &userRole); err != nil {
		panic(err)
	}
	return
}
