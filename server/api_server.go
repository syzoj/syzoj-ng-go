package server

import (
	"context"
	"net/http"
	"reflect"
	"sync"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"github.com/syzoj/syzoj-ng-go/model"
)

type ApiServer struct {
	s *Server

	debug       bool
	ctx         context.Context
	router      *mux.Router
	wsUpgrader  websocket.Upgrader
	wg          sync.WaitGroup
	wsConn      map[*websocket.Conn]struct{}
	wsConnMutex sync.Mutex
	cancelFunc  func()
}

type apiContextKey struct{}
type ApiContext struct {
	r   *http.Request
	w   http.ResponseWriter
	s   *ApiServer
	mut []*model.Mutation
}

var jsonMarshaler = jsonpb.Marshaler{OrigName: true}
var jsonUnmarshaler = jsonpb.Unmarshaler{}

type ApiConfig struct {
	Debug bool `json:"debug"`
}

func (s *Server) newApiServer(cfg *ApiConfig) *ApiServer {
	ApiServer := new(ApiServer)
	ApiServer.s = s
	if cfg.Debug {
		ApiServer.debug = true
	}
	router := mux.NewRouter()
	router.PathPrefix("/api").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Allow", "OPTIONS, GET, HEAD, POST")
		w.WriteHeader(200)
	})
	ApiServer.router = router
	ApiServer.ctx, ApiServer.cancelFunc = context.WithCancel(s.ctx)
	ApiServer.wg.Add(1)
	return ApiServer
}

func (s *Server) ApiServer() *ApiServer {
	return s.apiServer
}

func (s *ApiServer) close() {
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

func (s *ApiServer) Router() *mux.Router {
	return s.router
}

func (s *ApiServer) WrapHandler(h interface{}, checkToken bool) http.Handler {
	if s.debug {
		checkToken = false
	}
	val := reflect.ValueOf(h)
	if val.Kind() != reflect.Func {
		panic("wrapHandler: Invalid handler passed in")
	}
	t := val.Type()
	var reqType reflect.Type
	if t.NumIn() == 1 {
	} else if t.NumIn() == 2 {
		reqType = t.In(1)
		if !reqType.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
			panic("wrapHandler: Type of the second input argument does not implement proto.Message")
		}
	} else {
		panic("wrapHandler: Number of input arguments is neither 1 nor 2")
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
	respType := t.Out(0)
	if !respType.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
		panic("wrapHandler: Type of first output does not ipmlement proto.Message")
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := &ApiContext{r: r, w: w, s: s}
		ctx := s.ctx
		ctx, cancelFunc := context.WithCancel(ctx)
		defer cancelFunc()
		ctx = context.WithValue(ctx, apiContextKey{}, c)
		defer c.Send()
		if checkToken {
			token := r.Header.Get("X-CSRF-Token")
			if token != "1" {
				c.SendError(ErrCSRF)
				return
			}
		}
		var out []reflect.Value
		if reqType != nil {
			reqValue := reflect.New(reqType.Elem())
			err := c.ReadBody(reqValue.Interface().(proto.Message))
			if err != nil {
				c.SendError(err)
				return
			}
			out = val.Call([]reflect.Value{reflect.ValueOf(ctx), reqValue})
		} else {
			out = val.Call([]reflect.Value{reflect.ValueOf(ctx)})
		}
		if out[1].Interface() != nil {
			c.SendError(out[1].Interface().(error))
		} else {
			c.SendBody(out[0].Interface().(proto.Message))
		}
	})
}

func (s *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.wg.Add(1)
	defer s.wg.Done()
	if s.debug {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
	}
	s.router.ServeHTTP(w, r)
}

func GetApiContext(ctx context.Context) *ApiContext {
	return ctx.Value(apiContextKey{}).(*ApiContext)
}

func (c *ApiContext) ReadBody(val proto.Message) error {
	return jsonUnmarshaler.Unmarshal(c.r.Body, val)
}

func (c *ApiContext) Mutate(path string, method string, val proto.Message) {
	m, err := ptypes.MarshalAny(val)
	if err != nil {
		log.WithError(err).Error("Failed to marshal message into any")
		return
	}
	mutation := &model.Mutation{
		Path:   proto.String(path),
		Method: proto.String(method),
		Value:  m,
	}
	c.mut = append(c.mut, mutation)
}

func (c *ApiContext) SendBody(val proto.Message) {
	c.Mutate("", "setBody", val)
}

func (c *ApiContext) SendError(err error) {
	c.Mutate("", "setError", &model.Error{Error: proto.String(err.Error())})
}

func (c *ApiContext) Send() {
	resp := &model.Response{Mutations: c.mut}
	if err := jsonMarshaler.Marshal(c.w, resp); err != nil {
		log.WithError(err).Error("Failed to send response")
	}
}

func (c *ApiContext) UpgradeWebSocket() (*websocket.Conn, error) {
	return c.s.wsUpgrader.Upgrade(c.w, c.r, nil)
}
