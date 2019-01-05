package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/dgraph-io/dgo"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

type ApiServer struct {
	router       *mux.Router
	dgraph       *dgo.Dgraph
}
type ApiContext struct {
	w            http.ResponseWriter
	r            *http.Request
	sessResponse *SessionResponse
}
type ApiResponse struct {
	Error   string           `json:"error,omitempty"`
	Data    interface{}      `json:"data,omitempty"`
	Session *SessionResponse `json:"session,omitempty"`
}
type SessionResponse struct {
	UserName string `json:"user_name"`
	LoggedIn bool   `json:"logged_in"`
}

func CreateApiServer(dgraph *dgo.Dgraph) (*ApiServer, error) {
	srv := &ApiServer{
		dgraph:       dgraph,
	}
	srv.setupRoutes()
	return srv, nil
}

func (srv *ApiServer) setupRoutes() {
	router := mux.NewRouter()
	router.Handle("/api/auth/register", srv.wrapHandlerWithBody(srv.HandleAuthRegister)).Methods("POST")
	router.Handle("/api/auth/login", srv.wrapHandlerWithBody(srv.HandleAuthLogin)).Methods("POST")
	router.Handle("/api/auth/logout", srv.wrapHandlerWithBody(srv.HandleAuthLogout)).Methods("POST")
	router.Handle("/api/problem/create", srv.wrapHandlerWithBody(srv.HandleProblemCreate)).Methods("POST")
	router.Handle("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/view", srv.wrapHandler(srv.HandleProblemView)).Methods("GET")
	/*
		router.Handle("/api/problemset/create", srv.wrapHandlerWithBody(srv.HandleCreateProblemset)).Methods("POST")
		router.Handle("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/add", srv.wrapHandlerWithBody(srv.HandleProblemsetAdd)).Methods("POST")
		router.Handle("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/list", srv.wrapHandlerWithBody(srv.HandleProblemsetList)).Methods("GET")
		router.Handle("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/view", srv.wrapHandlerWithBody(srv.HandleProblemsetView)).Methods("GET")
		router.Handle("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/submit", srv.wrapHandlerWithBody(srv.HandleProblemsetSubmit)).Methods("POST")
		router.Handle("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/reset-token", srv.wrapHandlerWithBody(srv.HandleResetProblemToken)).Methods("POST")
		router.Handle("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/update", srv.wrapHandlerWithBody(srv.HandleProblemUpdate)).Methods("POST")
		router.Handle("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/change-title", srv.wrapHandlerWithBody(srv.HandleProblemChangeTitle)).Methods("POST")
	*/
	srv.router = router
}

func (srv *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *ApiServer) wrapHandlerWithBody(handler func(*ApiContext) ApiError) http.Handler {
	return srv.wrapHandler(func(c *ApiContext) ApiError {
		var err error
		token := c.r.Header.Get("X-CSRF-Token")
		var cookie *http.Cookie
		if cookie, err = c.r.Cookie("CSRF"); err != nil {
			return ErrCSRF
		}
		if cookie.Value != token {
			return ErrCSRF
		}
		return handler(c)
	})
}

func (srv *ApiServer) wrapHandler(handler func(*ApiContext) ApiError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ApiContext{w: w, r: r}
		var err error
		defer func() {
			if err != nil {
				writeError(c, err.(ApiError))
				return
			}
		}()
		if _, err = c.r.Cookie("CSRF"); err != nil {
			var token [16]byte
			if _, err = rand.Read(token[:]); err != nil {
				err = internalServerError(err)
				return
			}
			http.SetCookie(c.w, &http.Cookie{
				Name:  "CSRF",
				Value: hex.EncodeToString(token[:]),
				Path:  "/",
			})
		}
		err = handler(c)
	})
}

func writeError(c *ApiContext, err ApiError) {
	if ierr, ok := err.(internalServerErrorType); ok {
		log.Errorf("Error handling request %s: %s", c.r.URL, ierr.Err)
	} else {
		log.Infof("Failed to handle request %s: %s", c.r.URL, err)
	}
	var err2 error
	defer func() {
		if err2 != nil {
			log.WithField("error", err2).Warning("Failed to write error")
		}
	}()
	c.w.WriteHeader(err.Code())
	encoder := json.NewEncoder(c.w)
	err2 = encoder.Encode(ApiResponse{Error: err.Error(), Session: c.sessResponse})
}

func writeResponse(c *ApiContext, data interface{}) {
	encoder := json.NewEncoder(c.w)
	var err error
	defer func() {
		if err != nil {
			log.WithField("error", err).Warning("Failed to write response")
		}
	}()
	err = encoder.Encode(ApiResponse{Data: data, Session: c.sessResponse})
}
