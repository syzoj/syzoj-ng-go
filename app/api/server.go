package api

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"

	"github.com/dgraph-io/dgo"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fastjson"

	"github.com/syzoj/syzoj-ng-go/app/judge"
)

var log = logrus.StandardLogger()

type ApiServer struct {
	router       *mux.Router
	dgraph       *dgo.Dgraph
	judgeService judge.Service
	parserPool   fastjson.ParserPool
	config       Config
}
type Config struct {
	DebugToken string `json:"debug_token"`
}

// Creates an API server.
func CreateApiServer(dgraph *dgo.Dgraph, judgeService judge.Service, config Config) (*ApiServer, error) {
	srv := &ApiServer{
		dgraph:       dgraph,
		judgeService: judgeService,
		config:       config,
	}
	srv.setupRoutes()
	return srv, nil
}

func (srv *ApiServer) setupRoutes() {
	router := mux.NewRouter()
	router.Handle("/api/register", srv.wrapHandler(Handle_Register)).Methods("POST")
	router.Handle("/api/login", srv.wrapHandler(Handle_Login)).Methods("POST")
	router.Handle("/api/nav/logout", srv.wrapHandler(Handle_Nav_Logout)).Methods("POST")
	router.Handle("/api/problem-db/my", srv.wrapHandler(Handle_ProblemDb_My)).Methods("GET")
	router.Handle("/api/problem-db/new", srv.wrapHandler(Handle_ProblemDb_New)).Methods("POST")
	router.Handle("/api/problem-db/view/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", srv.wrapHandler(Handle_ProblemDb_View)).Methods("GET")
	router.Handle("/api/problem-db/view/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/publicize", srv.wrapHandler(Handle_ProblemDb_View_Publicize)).Methods("POST")
	router.Handle("/api/problem-db/view/{problem_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/submit", srv.wrapHandler(Handle_ProblemDb_View_Submit)).Methods("POST")
	router.Handle("/api/problems", srv.wrapHandler(Handle_Problems)).Methods("GET")
	router.Handle("/api/problem/{name}", srv.wrapHandler(Handle_Problem)).Methods("GET")
	router.Handle("/api/submission/view/{submission_id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}", srv.wrapHandler(Handle_Submission_View)).Methods("GET")
	router.Handle("/api/submission/my", srv.wrapHandler(Handle_Submission_My)).Methods("GET")
	debugRouter := mux.NewRouter()
	if srv.config.DebugToken != "" {
		router.PathPrefix("/api/debug/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Debug-Token")
			if token != "" && subtle.ConstantTimeCompare([]byte(token), []byte(srv.config.DebugToken)) == 1 {
				debugRouter.ServeHTTP(w, r)
			} else {
				http.Error(w, "Token mismatch", 403)
			}
		})
	}
	srv.router = router
}

func (srv *ApiServer) wrapHandler(h func(*ApiContext) ApiError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ApiContext{
			res: w,
			req: r,
			srv: srv,
		}
		var err error
		if _, err = c.req.Cookie("CSRF"); err != nil {
			var token [16]byte
			if _, err = rand.Read(token[:]); err != nil {
				panic(err)
			}
			c.SetCookie(&http.Cookie{
				Name:  "CSRF",
				Value: hex.EncodeToString(token[:]),
				Path:  "/",
			})
		}
		token := c.GetHeader("X-CSRF-Token")
		cookie_token := c.GetCookie("CSRF")
		if cookie_token == "" || cookie_token != token {
			c.SendError(ErrCSRF)
			return
		}
		apiErr := h(c)
		if apiErr != nil {
			c.SendError(apiErr)
		}
	})
}

// Implements http.Handler interface. Serves HTTP requests.
func (srv *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}
