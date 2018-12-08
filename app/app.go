package app

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/syzoj/syzoj-ng-go/app/api"
	"github.com/syzoj/syzoj-ng-go/app/auth"
	"github.com/syzoj/syzoj-ng-go/app/git"
	judge_traditional "github.com/syzoj/syzoj-ng-go/app/judge/traditional"
	"github.com/syzoj/syzoj-ng-go/app/problemset"
	"github.com/syzoj/syzoj-ng-go/app/session"
)

type App struct {
	db      *sql.DB
	levelDB *leveldb.DB

	sessService       session.SessionService
	authService       auth.AuthService
	problemsetService problemset.ProblemsetService
	httpServer        *http.Server
	router            *mux.Router
	gitServer         *git.GitServer
	apiServer         *api.ApiServer

	traditionalJudgeService judge_traditional.TraditionalJudgeService
}

var MissingDependencyError = errors.New("Missing dependency")
var DoubleSetupError = errors.New("Attempting to setup a service twice")

func (app *App) SetupDB(conn string) error {
	if app.db != nil {
		return DoubleSetupError
	}
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}
	app.db = db
	return nil
}

func (app *App) SetupLevelDB(path string) (err error) {
	if app.levelDB != nil {
		return DoubleSetupError
	}
	app.levelDB, err = leveldb.OpenFile(path, nil)
	return
}

func (app *App) SetupSessionService() (err error) {
	if app.levelDB == nil {
		return MissingDependencyError
	}
	if app.sessService != nil {
		return DoubleSetupError
	}
	app.sessService, err = session.NewLevelDBSessionService(app.levelDB)
	return
}

func (app *App) SetupAuthService() (err error) {
	if app.levelDB == nil {
		return MissingDependencyError
	}
	if app.authService != nil {
		return DoubleSetupError
	}
	app.authService, err = auth.NewLevelDBAuthService(app.levelDB)
	return
}

func (app *App) SetupProblemsetService() (err error) {
	if app.levelDB == nil || app.traditionalJudgeService == nil {
		return MissingDependencyError
	}
	if app.problemsetService != nil {
		return DoubleSetupError
	}
	app.problemsetService, err = problemset.NewProblemsetService(app.levelDB, app.traditionalJudgeService)
	app.traditionalJudgeService.RegisterProblemsetService(app.problemsetService)
	return
}

func (app *App) SetupTraditionalJudgeService() (err error) {
	if app.traditionalJudgeService != nil {
		return DoubleSetupError
	}
	app.traditionalJudgeService, err = judge_traditional.NewTraditionalJudgeService()
	return
}

func (app *App) SetupGitServer(gitPath string) error {
	if app.db == nil {
		return MissingDependencyError
	}
	if app.gitServer != nil {
		return DoubleSetupError
	}
	server, err := git.CreateGitServer(app.db, gitPath)
	if err != nil {
		return err
	}
	app.gitServer = server
	return nil
}

func (app *App) SetupApiServer() error {
	if app.authService == nil || app.sessService == nil || app.problemsetService == nil {
		return MissingDependencyError
	}
	if app.apiServer != nil {
		return DoubleSetupError
	}
	server, err := api.CreateApiServer(app.sessService, app.authService, app.problemsetService)
	if err != nil {
		return err
	}

	app.apiServer = server
	return nil
}

func (app *App) SetupHttpServer(addr string) error {
	if app.httpServer != nil {
		return DoubleSetupError
	}
	app.router = mux.NewRouter()
	app.httpServer = &http.Server{
		Addr:         addr,
		Handler:      app.router,
		WriteTimeout: time.Second * 10,
	}
	return nil
}

func (app *App) AddGitServer() error {
	if app.router == nil {
		return MissingDependencyError
	}
	gitRouter := mux.NewRouter()
	gitRouter.HandleFunc("/git/{git-id}/HEAD", app.gitServer.Handle).Methods("GET")
	gitRouter.HandleFunc("/git/{git-id}/info/refs", app.gitServer.Handle).Methods("GET")
	gitRouter.HandleFunc("/git/{git-id}/objects/info/alternates", app.gitServer.Handle).Methods("GET")
	gitRouter.HandleFunc("/git/{git-id}/objects/info/http-alternates", app.gitServer.Handle).Methods("GET")
	gitRouter.HandleFunc("/git/{git-id}/objects/info/packs", app.gitServer.Handle).Methods("GET")
	gitRouter.HandleFunc("/git/{git-id}/objects/{object-id-prefix:[0-9a-f]{2}}/{object-id-suffix:[0-9a-f]{38}}", app.gitServer.Handle).Methods("GET")
	gitRouter.HandleFunc("/git/{git-id}/objects/pack/pack-{pack-id:[0-9a-f]{40}}.pack", app.gitServer.Handle).Methods("GET")
	gitRouter.HandleFunc("/git/{git-id}/objects/pack/pack-{pack-id:[0-9a-f]{40}}.idx", app.gitServer.Handle).Methods("GET")
	gitRouter.HandleFunc("/git/{git-id}/git-upload-pack", app.gitServer.Handle).Methods("POST")
	gitRouter.HandleFunc("/git/{git-id}/git-receive-pack", app.gitServer.Handle).Methods("POST")
	app.router.PathPrefix("/git/{git-id:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/").Handler(gitRouter)
	return nil
}

func (app *App) AddApiServer() error {
	if app.router == nil {
		return MissingDependencyError
	}
	app.router.PathPrefix("/api").Handler(app.apiServer)
	return nil
}

func (app *App) runWebServer() {
	log.Println("Starting web server at", app.httpServer.Addr)
	err := app.httpServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Println("Web server failed unexpectedly: ", err)
	}
}

func (app *App) Run() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go app.runWebServer()
	<-sigChan
	app.httpServer.Shutdown(context.Background())
	if app.sessService != nil {
		log.Println("Shutting down session service")
		app.sessService.Close()
	}
	if app.authService != nil {
		log.Println("Shutting down auth service")
		app.authService.Close()
	}
	if app.problemsetService != nil {
		log.Println("Shutting down problemset service")
		app.problemsetService.Close()
	}
	if app.traditionalJudgeService != nil {
		log.Println("Shutting down traditional judge service")
		app.traditionalJudgeService.Close()
	}
	if app.levelDB != nil {
		log.Println("Shutting down LevelDB")
		app.levelDB.Close()
	}
	log.Println("Server shut down")
}
