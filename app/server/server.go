package server

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

type Server struct {
	db         *sql.DB
	ctx        context.Context
	cancelFunc func()

	apiServer *apiServer
}

func NewServer(db *sql.DB) *Server {
	server := new(Server)
	server.db = db
	server.ctx, server.cancelFunc = context.WithCancel(context.Background())
	server.apiServer = server.newApiServer()
	return server
}

func (s *Server) Close() error {
	s.apiServer.close()
	return nil
}
