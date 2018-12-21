package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/syzoj/syzoj-ng-go/app/auth"
	"github.com/syzoj/syzoj-ng-go/app/judge"
	"github.com/syzoj/syzoj-ng-go/app/problemset"
	"github.com/syzoj/syzoj-ng-go/app/session"
)

type ApiServer struct {
	router            *mux.Router
	sessService       session.Service
	authService       auth.Service
	problemsetService problemset.Service
	judgeService      judge.Service
}

var defaultUserId = uuid.MustParse("00000000-0000-0000-0000-000000000000")

func CreateApiServer(sessService session.Service, authService auth.Service, problemsetService problemset.Service, judgeService judge.Service) (*ApiServer, error) {
	srv := &ApiServer{
		sessService:       sessService,
		authService:       authService,
		problemsetService: problemsetService,
		judgeService:      judgeService,
	}
	srv.setupRoutes()
	return srv, nil
}

func (srv *ApiServer) setupRoutes() {
	router := mux.NewRouter()
	router.Handle("/api/auth/register", srv.wrapHandlerWithSession(srv.HandleAuthRegister)).Methods("POST")
	router.Handle("/api/auth/login", srv.wrapHandlerWithSession(srv.HandleAuthLogin)).Methods("POST")
	router.Handle("/api/auth/logout", srv.wrapHandlerWithSession(srv.HandleAuthLogout)).Methods("POST")
	router.Handle("/api/problemset/create", srv.wrapHandlerWithSession(srv.HandleCreateProblemset)).Methods("POST")
	router.Handle("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/add", srv.wrapHandlerWithSession(srv.HandleProblemsetAdd)).Methods("POST")
	router.Handle("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/list", srv.wrapHandlerWithSession(srv.HandleProblemsetList)).Methods("GET")
	router.Handle("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/view", srv.wrapHandlerWithSession(srv.HandleProblemsetView)).Methods("GET")
	router.Handle("/api/problemset/{problemset_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/submit", srv.wrapHandlerWithSession(srv.HandleProblemsetSubmit)).Methods("POST")
	router.Handle("/api/problem/create", srv.wrapHandlerWithSession(srv.HandleProblemCreate)).Methods("POST")
	router.Handle("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/view", srv.wrapHandlerWithSession(srv.HandleProblemView)).Methods("GET")
	router.Handle("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/reset-token", srv.wrapHandlerWithSession(srv.HandleResetProblemToken)).Methods("POST")
	router.Handle("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/update", srv.wrapHandlerWithSession(srv.HandleProblemUpdate)).Methods("POST")
	router.Handle("/api/problem/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/change-title", srv.wrapHandlerWithSession(srv.HandleProblemChangeTitle)).Methods("POST")
	srv.router = router
}

func (srv *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *ApiServer) ensureSession(w http.ResponseWriter, r *http.Request) (uuid.UUID, *session.Session, ApiError) {
	var sessId uuid.UUID
	if cookie, err := r.Cookie("SYZOJSESSION"); err == nil {
		sessId, _ = uuid.Parse(cookie.Value)
	}
	if sess, err := srv.sessService.GetSession(sessId); err != nil {
		if err != session.ErrSessionNotFound {
			return sessId, sess, internalServerError(err)
		}
		if sessId, sess, err := srv.sessService.NewSession(); err == nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "SYZOJSESSION",
				Value:    sessId.String(),
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Now().Add(time.Hour * 24 * 30),
			})
			return sessId, sess, nil
		} else {
			return sessId, sess, internalServerError(err)
		}
	} else {
		return sessId, sess, nil
	}
}

func (srv *ApiServer) updateSession(sessId uuid.UUID, sess *session.Session) ApiError {
	var err = srv.sessService.UpdateSession(sessId, sess)
	if err != nil {
		return internalServerError(err)
	}
	return nil
}

func (srv *ApiServer) wrapHandler(handler func(http.ResponseWriter, *http.Request) ApiError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err = handler(w, r)
		if err != nil {
			writeError(w, r, err, nil)
			return
		}
	})
}

func (srv *ApiServer) wrapHandlerWithSession(handler func(http.ResponseWriter, *http.Request, uuid.UUID, *session.Session) ApiError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessId, sess, err := srv.ensureSession(w, r)
		if err != nil {
			writeError(w, r, err, nil)
			return
		}
		if err = handler(w, r, sessId, sess); err != nil {
			writeError(w, r, err, sess)
			return
		}
	})
}

func requireLogin(sess *session.Session) ApiError {
	if sess.AuthUserId == defaultUserId {
		return ErrNotLoggedIn
	} else {
		return nil
	}
}
