package api

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/gorilla/mux"

	"github.com/syzoj/syzoj-ng-go/app/auth"
	"github.com/syzoj/syzoj-ng-go/app/session"
)

type ApiServer struct {
	router      *mux.Router
	sessService session.SessionService
	authService auth.AuthService
}

func CreateApiServer(sessService session.SessionService, authService auth.AuthService) (*ApiServer, error) {
	srv := &ApiServer{
		sessService: sessService,
		authService: authService,
	}
	srv.setupRoutes()
	return srv, nil
}

func (srv *ApiServer) setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/api/auth/register", srv.HandleAuthRegister)
	router.HandleFunc("/api/auth/login", srv.HandleAuthLogin)
	srv.router = router
}

func (srv *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *ApiServer) ensureSession(r *http.Request) (uuid.UUID, *session.Session, error) {
	var sessId uuid.UUID
	if cookie, err := r.Cookie("SYZOJSESSION"); err == nil {
		sessId, _ = uuid.Parse(cookie.Value)
	}
	if sess, err := srv.sessService.GetSession(sessId); err != nil {
		if sessId, sess, err := srv.sessService.NewSession(); err != nil {
			return sessId, sess, err
		} else {
			return sessId, sess, err
		}
	} else {
		return sessId, sess, nil
	}
}
