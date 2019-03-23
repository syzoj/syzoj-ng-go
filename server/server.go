package server

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/database"
)

var log = logrus.StandardLogger()

type Server struct {
	db         *database.Database
	ctx        context.Context
	cancelFunc func()

	apiServer *ApiServer
}

type serverKey struct{}

type ServerConfig struct {
	API ApiConfig `json:"api"`
}

func NewServer(db *database.Database, cfg *ServerConfig) *Server {
	server := new(Server)
	server.db = db
	ctx := context.Background()
	ctx = context.WithValue(ctx, serverKey{}, server)
	server.ctx, server.cancelFunc = context.WithCancel(ctx)
	server.apiServer = server.newApiServer(&cfg.API)
	return server
}

func GetServer(ctx context.Context) *Server {
	return ctx.Value(serverKey{}).(*Server)
}

func (s *Server) GetDB() *database.Database {
	return s.db
}

func (s *Server) Close() error {
	s.apiServer.close()
	return nil
}
