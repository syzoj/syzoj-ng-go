package api

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"

	"github.com/syzoj/syzoj-ng-go/app/judge"
	model_group "github.com/syzoj/syzoj-ng-go/app/model/group"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type ApiServer struct {
	db           *sql.DB
	judgeService judge.JudgeServiceProvider
	redis        *redis.Client
	router       *mux.Router
}
type ApiContext struct {
	w    http.ResponseWriter
	r    *http.Request
	s    *ApiServer
	sess *Session
	resp interface{}
	code int

	groupName      string
	groupId        util.UUID
	problemsetName string
	problemsetId   util.UUID
	problemName    string
	problemId      util.UUID

	groupPolicy   model_group.GroupPolicy
	groupUserRole model_group.GroupUserRole

	tx        *sql.Tx
	txSuccess bool
}
type ApiHandler func(*ApiContext) ApiResponse
type ApiResponse interface {
    Execute(cxt *ApiContext)
}
type Session struct {
	SessionId  string
	AuthUserId util.UUID
	willSave   bool
}

type ErrorResponse struct {
	Error string `json:"error"`
}
type SuccessResponse struct {
	Data interface{} `json:"data"`
}
type SuccessResponseType struct {
    Value interface{}
}

func CreateApiServer(db *sql.DB, redis *redis.Client, judgeService judge.JudgeServiceProvider) (*ApiServer, error) {
	router := mux.NewRouter()
	srv := &ApiServer{
		db:           db,
		judgeService: judgeService,
		redis:        redis,
		router:       router,
	}
	router.Handle("/api/auth/register", srv.ApiHandler(HandleAuthRegister)).Methods("POST")
	router.Handle("/api/auth/login", srv.ApiHandler(HandleAuthLogin)).Methods("POST")
	router.Handle("/api/group/create", srv.ApiHandler(HandleGroupCreate)).Methods("POST")
	router.Handle("/api/user/info", srv.ApiHandler(HandleUserInfo)).Methods("GET")
	router.Handle("/api/group/problemset/create", srv.ApiHandler(HandleProblemsetCreate)).Methods("POST")
	/*
	   router.Handle("/api/group/problemset/problem/create", srv.ApiHandler(HandleProblemCreate)).Methods("POST")
	*/
	return srv, nil
}

func (srv *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func marshalJson(data interface{}) []byte {
	result, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return result
}

func Success(resp interface{}) ApiResponse {
    return SuccessResponseType{Value: resp}
}
func (r SuccessResponseType) Execute(cxt *ApiContext) {
    cxt.code = 200
    cxt.resp = SuccessResponse{r.Value}
}
func (cxt *ApiContext) Complete() {
	cxt.w.WriteHeader(cxt.code)
	if data, err := json.Marshal(cxt.resp); err != nil {
		panic(err)
	} else {
		cxt.w.Write(data)
	}
}
func (cxt *ApiContext) ReadBody(body interface{}) *ApiError {
	decoder := json.NewDecoder(cxt.r.Body)
	if err := decoder.Decode(body); err != nil {
		log.Println("Bad request:", err)
		return BadRequestError
	}
	return nil
}

func (sess *Session) Save() {
	sess.willSave = true
}
func (sess *Session) IsLoggedIn() bool {
	return sess.AuthUserId != (util.UUID{})
}

func (srv *ApiServer) ApiHandler(h ApiHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cxt := &ApiContext{
			r:    r,
			w:    w,
			s:    srv,
			code: 200,
		}
		defer cxt.Complete()
		defer func() {
			err := recover()
			if err != nil {
				log.Println("Error handling http request:", err)
				debug.PrintStack()
				cxt.code = 500
				cxt.resp = ErrorResponse{Error: "Internal server error"}
			}
		}()

		sess := &Session{}
		if cookie, err := cxt.r.Cookie("SYZOJSESSION"); cookie != nil && err == nil {
			sessId := cookie.Value
			sess_map, err := cxt.s.redis.HGetAll(fmt.Sprintf("sess:%s", sessId)).Result()
			if err != redis.Nil && len(sess_map) != 0 {
				sess.SessionId = sessId
				if val, ok := sess_map["user-id"]; ok {
					if userId, err := util.UUIDFromBytes([]byte(val)); err != nil {
						panic(InvalidAuthUserIdError)
					} else {
						sess.AuthUserId = userId
					}
				}
			}
		}
		cxt.sess = sess
		defer func() {
			if cxt.sess.willSave {
				var setExpire = false
				if cxt.sess.SessionId == "" {
					var buf [32]byte
					_, err := rand.Read(buf[:])
					if err != nil {
						panic(err)
					}

					sessId := make([]byte, 64)
					hex.Encode(sessId, buf[:])
					cxt.sess.SessionId = string(sessId)
					http.SetCookie(w, &http.Cookie{
						Name:     "SYZOJSESSION",
						Value:    sess.SessionId,
						Expires:  time.Now().Add(24 * time.Hour),
						HttpOnly: true,
					})
					setExpire = true
				}

				key := "sess:" + cxt.sess.SessionId
				if _, err := srv.redis.Del(key).Result(); err != nil {
					panic(err)
				}
				rmap := make(map[string]interface{})
				if cxt.sess.IsLoggedIn() {
					rmap["user-id"] = cxt.sess.AuthUserId
				}
				if _, err := srv.redis.HMSet(key, rmap).Result(); err != nil {
					panic(err)
				}

				if setExpire {
					if _, err := srv.redis.Expire(key, 24*time.Hour).Result(); err != nil {
						panic(err)
					}
				}
			}
		}()

		defer func() {
			if cxt.tx != nil {
				if cxt.txSuccess {
					if err := cxt.tx.Commit(); err != nil {
						panic(err)
					}
				} else {
					if err := cxt.tx.Rollback(); err != nil {
						panic(err)
					}
				}
			}
		}()

		if resp := h(cxt); resp != nil {
            resp.Execute(cxt)
		}
	})
}

func (srv *ApiServer) HandleCatchAll(ctx *ApiContext) ApiResponse {
	return ApiEndpointNotFoundError
}

func UseTx(cxt *ApiContext) ApiResponse {
	if cxt.tx != nil {
		panic("UseTx called twice")
	}
	if tx, err := cxt.s.db.Begin(); err != nil {
		panic(err)
	} else {
		cxt.tx = tx
	}
	return nil
}

func DoneTx(cxt *ApiContext) ApiResponse {
	cxt.txSuccess = true
	return nil
}
