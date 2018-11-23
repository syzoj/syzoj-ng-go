package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (srv *ApiServer) HandleProblemCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars
}
