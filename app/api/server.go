package api

import (
	"context"
	"crypto/subtle"
	"net/http"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"

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

	streamSender   map[string]*StreamSender
	streamReceiver map[string]*StreamReceiver
	streamLock     sync.Mutex
}
type Config struct {
	DebugToken string `json:"debug_token"`
}

var jsonMarshaler = jsonpb.Marshaler{OrigName: true}
var jsonUnmarshaler = jsonpb.Unmarshaler{}

// Creates an API server.
func CreateApiServer(mongodb *mongo.Client, c *core.Core, config Config) (*ApiServer, error) {
	srv := &ApiServer{
		mongodb:        mongodb.Database("syzoj"),
		c:              c,
		config:         config,
		wsConn:         make(map[*websocket.Conn]struct{}),
		streamSender:   make(map[string]*StreamSender),
		streamReceiver: make(map[string]*StreamReceiver),
	}
	srv.wg.Add(1)
	srv.ctx, srv.cancelFunc = context.WithCancel(context.Background())
	srv.setupRoutes()
	return srv, nil
}

func (srv *ApiServer) setupRoutes() {
	router := mux.NewRouter()
	router.Handle("/api/stream/{token}", http.HandlerFunc(srv.HandleStream)).Methods("POST", "PUT")
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

func (srv *ApiServer) wrapHandler(h func(context.Context, *ApiContext) ApiError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ApiContext{
			res:    w,
			req:    r,
			Server: srv,
		}
		token := c.GetHeader("X-CSRF-Token")
		if token != "1" {
			c.SendError(ErrCSRF)
			return
		}
		apiErr := h(r.Context(), c)
		if apiErr != nil {
			c.SendError(apiErr)
		}
	})
}

func (srv *ApiServer) wrapHandlerNoToken(h func(context.Context, *ApiContext) ApiError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ApiContext{
			res:    w,
			req:    r,
			Server: srv,
		}
		apiErr := h(r.Context(), c)
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
	srv.streamLock.Lock()
	for _, sender := range srv.streamSender {
		close(sender.closeChan)
	}
	srv.streamLock.Unlock()
	srv.wg.Done()
	srv.wg.Wait()
}
