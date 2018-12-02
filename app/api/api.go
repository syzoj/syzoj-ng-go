package api

import (
	"database/sql"
	"net/http"

	"github.com/syzoj/syzoj-ng-go/app/lock"

	"github.com/gorilla/mux"

	"github.com/go-redis/redis"
	"github.com/syzoj/syzoj-ng-go/app/judge"
)

type ApiServer struct {
	db           *sql.DB
	redis        *redis.Client
	judgeService judge.JudgeServiceProvider
	router       *mux.Router
	lockManager  lock.LockManager
}

func CreateApiServer(db *sql.DB, redis *redis.Client, judgeService judge.JudgeServiceProvider, lockManager lock.LockManager) (*ApiServer, error) {
	srv := &ApiServer{
		db:           db,
		redis:        redis,
		judgeService: judgeService,
		lockManager:  lockManager,
	}
	srv.setupRoutes()
	return srv, nil
}

func (srv *ApiServer) setupRoutes() {
	router := mux.NewRouter()
	router.HandleFunc("/api/auth/login", srv.HandleAuthLogin)
	//router.HandleFunc("/api/auth/register", srv.HandleAuthRegister)
	srv.router = router
}

func (srv *ApiServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}
