package api

import (
	"github.com/syzoj/syzoj-ng-go/app/problemset"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gorilla/mux"

	"github.com/syzoj/syzoj-ng-go/app/auth"
	"github.com/syzoj/syzoj-ng-go/app/session"
)

type ApiServer struct {
	router      *mux.Router
	sessService session.SessionService
	authService auth.AuthService
	problemsetService problemset.ProblemsetService
}

var defaultUserId = uuid.MustParse("00000000-0000-0000-0000-000000000000")

func CreateApiServer(sessService session.SessionService, authService auth.AuthService, problemsetService problemset.ProblemsetService) (*ApiServer, error) {
	srv := &ApiServer{
		sessService: sessService,
		authService: authService,
		problemsetService: problemsetService,
	}
	srv.setupRoutes()
	return srv, nil
}

func (srv *ApiServer) setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/api/auth/register", srv.HandleAuthRegister).Methods("POST")
	router.HandleFunc("/api/auth/login", srv.HandleAuthLogin).Methods("POST")
	router.HandleFunc("/api/problemset/regular/create", srv.HandleRegularProblemsetCreate).Methods("POST")
	srv.router = router
}

func (srv *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *ApiServer) ensureSession(w http.ResponseWriter, r *http.Request) (uuid.UUID, *session.Session, error) {
	var sessId uuid.UUID
	if cookie, err := r.Cookie("SYZOJSESSION"); err == nil {
		sessId, _ = uuid.Parse(cookie.Value)
	}
	if sess, err := srv.sessService.GetSession(sessId); err != nil {
		if sessId, sess, err := srv.sessService.NewSession(); err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "SYZOJSESSION",
				Value:    sessId.String(),
				HttpOnly: true,
				Expires:  time.Now().Add(time.Hour * 24),
			})
			return sessId, sess, err
		} else {
			return sessId, sess, err
		}
	} else {
		return sessId, sess, nil
	}
}
