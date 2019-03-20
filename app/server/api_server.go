package server

import (
	"context"
	"net/http"
	"sync"
    "reflect"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

    "github.com/syzoj/syzoj-ng-go/app/model"
)

type apiServer struct {
    s *Server

    ctx context.Context
	router      *mux.Router
	wsUpgrader  websocket.Upgrader
	wg          sync.WaitGroup
	wsConn      map[*websocket.Conn]struct{}
	wsConnMutex sync.Mutex
	cancelFunc  func()
}

type apiContextKey struct{}
type apiContext struct {
    w http.ResponseWriter
    r *http.Request
}

var jsonMarshaler = jsonpb.Marshaler{OrigName: true}
var jsonUnmarshaler = jsonpb.Unmarshaler{}

func (s *Server) newApiServer() *apiServer {
    apiServer := new(apiServer)
    router := mux.NewRouter()
    apiServer.router = router
    apiServer.setupRoutes()
    apiServer.ctx, apiServer.cancelFunc = context.WithCancel(s.ctx)
    apiServer.wg.Add(1)
    return apiServer
}

func (s *Server) ApiServer() http.Handler {
    return s.apiServer
}

func (s *apiServer) close() {
    s.cancelFunc()
	s.wsConnMutex.Lock()
	for conn := range s.wsConn {
		conn.Close()
	}
	s.wsConn = nil
	s.wsConnMutex.Unlock()
    s.wg.Done()
    s.wg.Wait()
}

func (s *apiServer) setupRoutes() {
    s.router.Path("/api/login").Methods("POST").Handler(s.wrapHandler(s.Handle_Login, true))
}

func (s *apiServer) wrapHandler(h interface{}, checkToken bool) http.Handler {
    val := reflect.ValueOf(h)
    if val.Kind() != reflect.Func {
        panic("wrapHandler: Invalid handler passed in")
    }
    t := val.Type()
    if t.NumIn() != 2 {
        panic("wrapHandler: Number of input arguments is not 2")
    }
    if t.NumOut() != 2 {
        panic("wrapHandler: Number of outputs is not 2")
    }
    if t.In(0) != reflect.TypeOf((*context.Context)(nil)).Elem() {
        panic("wrapHandler: Type of first argument is not context.Context")
    }
    if t.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
        panic("wrapHandler: Type of second output is not error")
    }
    reqType := t.In(1)
    respType := t.Out(0)
    if !reqType.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
        panic("wrapHandler: Type of second argument does not implement proto.Message")
    }
    if !respType.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
        panic("wrapHandler: Type of first output does not ipmlement proto.Message")
    }
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        c := &apiContext{r: r, w: w}
        ctx := context.WithValue(r.Context(), apiContextKey{}, c)
        if checkToken {
            token := r.Header.Get("X-CSRF-Token")
            if token != "1" {
                s.SendError(ctx, ErrCSRF)
                return
            }
        }
        reqValue := reflect.New(reqType.Elem())
        err := s.ReadBody(ctx, reqValue.Interface().(proto.Message))
        if err != nil {
            s.SendError(ctx, err)
            return
        }
        out := val.Call([]reflect.Value{reflect.ValueOf(ctx), reqValue})
        if out[1].Interface() != nil {
            s.SendError(ctx, out[1].Interface().(error))
        } else {
            s.SendBody(ctx, out[0].Interface().(proto.Message))
        }
    })
}

func (s *apiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.wg.Add(1)
	defer s.wg.Done()
	s.router.ServeHTTP(w, r)
}

func getApiContext(ctx context.Context) *apiContext {
    return ctx.Value(apiContextKey{}).(*apiContext)
}

func (s *apiServer) ReadBody(ctx context.Context, val proto.Message) error {
    c := getApiContext(ctx)
    return jsonUnmarshaler.Unmarshal(c.r.Body, val)
}

func (s *apiServer) SendBody(ctx context.Context, val proto.Message) {
    c := getApiContext(ctx)
    resp := new(model.Response)
    var err2 error
    resp.Data, err2 = ptypes.MarshalAny(val)
    if err2 != nil {
        log.WithError(err2).Error("Failed to send response")
        return
    }
    err2 = jsonMarshaler.Marshal(c.w, resp)
    if err2 != nil {
        log.WithError(err2).Error("Failed to send response")
        return
    }
}

func (s *apiServer) SendError(ctx context.Context, err error) {
    c := getApiContext(ctx)
    resp := new(model.Response)
    resp.Error = proto.String(err.Error())
    err2 := jsonMarshaler.Marshal(c.w, resp)
    if err2 != nil {
        log.WithError(err2).Error("Failed to send response")
        return
    }
}

func (s *apiServer) UpgradeWebSocket(ctx context.Context) (*websocket.Conn, error) {
    c := getApiContext(ctx)
    return s.wsUpgrader.Upgrade(c.w, c.r, nil)
}
