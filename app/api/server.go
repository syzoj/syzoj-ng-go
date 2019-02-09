package api

import (
	"context"
	"crypto/subtle"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/app/core"
)

var log = logrus.StandardLogger()

type ApiServer struct {
	router      *mux.Router
	mongodb     *mongo.Database
	c           *core.Core
	config      Config
	wsUpgrader  websocket.Upgrader
	wg          sync.WaitGroup
	ctx         context.Context
	wsConn      map[*websocket.Conn]struct{}
	wsConnMutex sync.Mutex
	cancelFunc  func()
}
type Config struct {
	DebugToken string `json:"debug_token"`
}

// Creates an API server.
func CreateApiServer(mongodb *mongo.Client, c *core.Core, config Config) (*ApiServer, error) {
	srv := &ApiServer{
		mongodb: mongodb.Database("syzoj"),
		c:       c,
		config:  config,
		wsConn:  make(map[*websocket.Conn]struct{}),
	}
	srv.wg.Add(1)
	srv.ctx, srv.cancelFunc = context.WithCancel(context.Background())
	srv.setupRoutes()
	return srv, nil
}

func (srv *ApiServer) setupRoutes() {
	router := mux.NewRouter()
	router.Handle("/api/register", srv.wrapHandler(Handle_Register)).Methods("POST")
	router.Handle("/api/login", srv.wrapHandler(Handle_Login)).Methods("POST")
	router.Handle("/api/nav/logout", srv.wrapHandler(Handle_Nav_Logout)).Methods("POST")
	router.Handle("/api/problem-db", srv.wrapHandler(Handle_ProblemDb)).Methods("GET")
	router.Handle("/api/problem-db/new", srv.wrapHandler(Handle_ProblemDb_New)).Methods("POST")
	router.Handle("/api/problem-db/view/{problem_id:[0-9A-Za-z\\-_]{16}}", srv.wrapHandler(Handle_ProblemDb_View)).Methods("GET")
	router.Handle("/api/problem-db/view/{problem_id:[0-9A-Za-z\\-_]{16}}/submit", srv.wrapHandler(Handle_ProblemDb_View_Submit)).Methods("POST")
	router.Handle("/api/problem-db/view/{problem_id:[0-9A-Za-z\\-_]{16}}/edit", srv.wrapHandler(Handle_ProblemDb_View_Edit)).Methods("POST")
	router.Handle("/api/contests", srv.wrapHandler(Handle_Contests)).Methods("GET")
	router.Handle("/api/contest-new", srv.wrapHandler(Handle_Contest_New)).Methods("POST")
	router.Handle("/api/contest/{contest_id:[0-9A-Za-z\\-_]{16}}/register", srv.wrapHandler(Handle_Contest_Register)).Methods("POST")
	router.Handle("/api/contest/{contest_id:[0-9A-Za-z\\-_]{16}}/index", srv.wrapHandler(Handle_Contest_Index)).Methods("GET")
	router.Handle("/api/contest/{contest_id:[0-9A-Za-z\\-_]{16}}/load", srv.wrapHandler(Handle_Contest_Load)).Methods("POST")
	router.Handle("/api/contest/{contest_id:[0-9A-Za-z\\-_]{16}}/unload", srv.wrapHandler(Handle_Contest_Unload)).Methods("POST")
	router.Handle("/api/contest/{contest_id:[0-9A-Za-z\\-_]{16}}/problem/{entry_name}", srv.wrapHandler(Handle_Contest_Problem)).Methods("GET")
	router.Handle("/api/contest/{contest_id:[0-9A-Za-z\\-_]{16}}/problem/{entry_name}/submit", srv.wrapHandler(Handle_Contest_Problem_Submit)).Methods("POST")
	router.Handle("/api/contest/{contest_id:[0-9A-Za-z\\-_]{16}}/status", srv.wrapHandlerNoToken(Handle_Contest_Status)).Methods("GET")
	router.Handle("/api/submissions", srv.wrapHandler(Handle_Submissions)).Methods("GET")
	router.Handle("/api/submission/view/{submission_id:[0-9A-Za-z\\-_]{16}}", srv.wrapHandler(Handle_Submission_View)).Methods("GET")
	router.Handle("/api/articles", srv.wrapHandler(Handle_Articles)).Methods("GET")
	router.Handle("/api/article/view/{article_id:[0-9A-Za-z\\-_]{16}}", srv.wrapHandler(Handle_Article_View)).Methods("GET")
	debugRouter := mux.NewRouter()
	debugRouter.Handle("/api/debug/submission/{submission_id:[0-9A-Za-z\\-_]{16}}/enqueue", srv.wrapHandlerNoToken(Handle_Debug_Submission_Enqueue)).Methods("POST")
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
		token := c.GetHeader("X-CSRF-Token")
		if token != "1" {
			c.SendError(ErrCSRF)
			return
		}
		apiErr := h(c)
		if apiErr != nil {
			c.SendError(apiErr)
		}
	})
}

func (srv *ApiServer) wrapHandlerNoToken(h func(*ApiContext) ApiError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ApiContext{
			res: w,
			req: r,
			srv: srv,
		}
		apiErr := h(c)
		if apiErr != nil {
			c.SendError(apiErr)
		}
	})
}

// Implements http.Handler interface. Serves HTTP requests.
func (srv *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.wg.Add(1)
	defer srv.wg.Done()
	srv.router.ServeHTTP(w, r)
}

func (srv *ApiServer) Close() {
	srv.cancelFunc()
	srv.wsConnMutex.Lock()
	for conn := range srv.wsConn {
		conn.Close()
	}
	srv.wsConn = nil // Cause errors if someone attempts to write to wsConn
	srv.wsConnMutex.Unlock()
	srv.wg.Done()
	srv.wg.Wait()
}
