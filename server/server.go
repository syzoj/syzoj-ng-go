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

	apiServer *apiServer
}

func NewServer(db *database.Database) *Server {
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
