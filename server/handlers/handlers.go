package handlers

import (
	"github.com/sirupsen/logrus"

	"github.com/syzoj/syzoj-ng-go/server"
)

var log = logrus.StandardLogger()

func RegisterHandlers(s *server.ApiServer) {
	router := s.Router()
	loginPage := router.PathPrefix("/api/login").Subrouter()
	loginPage.Path("").Methods("GET").Handler(s.WrapHandler(Get_Login, true))
	loginPage.Path("/page/login").Methods("POST").Handler(s.WrapHandler(Handle_Login, true))
	registerPage := router.PathPrefix("/api/register").Subrouter()
	registerPage.Path("").Methods("POST").Handler(s.WrapHandler(Handle_Register, true))
}
