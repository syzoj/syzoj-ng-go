package app

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syzoj/syzoj-ng-go/app/session"

	"github.com/gorilla/mux"

	"github.com/syzoj/syzoj-ng-go/app/api"
	"github.com/syzoj/syzoj-ng-go/app/auth"
	"github.com/syzoj/syzoj-ng-go/app/git"
)

type App struct {
	db      *sql.DB
	levelDB *leveldb.DB

	sessService session.SessionService
	authService auth.AuthService
	httpServer  *http.Server
	router      *mux.Router
	gitServer   *git.GitServer
	apiServer   *api.ApiServer
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
		return DoubleSetupError
	}
	if app.sessService != nil {
		return DoubleSetupError
	}
	app.sessService, err = session.NewLevelDBSessionService(app.levelDB)
	return
}

func (app *App) SetupAuthService() (err error) {
	if app.levelDB == nil {
		return DoubleSetupError
	}
	if app.authService != nil {
		return DoubleSetupError
	}
	app.authService, err = auth.NewLevelDBAuthService(app.levelDB)
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
	if app.authService == nil || app.sessService == nil {
		return MissingDependencyError
	}
	if app.apiServer != nil {
		return DoubleSetupError
	}
	server, err := api.CreateApiServer(app.sessService, app.authService)
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

func (app *App) runWebServer(done <-chan struct{}) {
	errChan := make(chan error)
	go func() {
		log.Println("Starting web server at", app.httpServer.Addr)
		errChan <- app.httpServer.ListenAndServe()
	}()
	select {
	case <-done:
		log.Println("Shutting down web server")
		app.httpServer.Shutdown(context.Background())
	case err := <-errChan:
		if err != http.ErrServerClosed {
			log.Println("Web server failed unexpectedly: ", err)
		}
	}
}

func (app *App) Run() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	doneChan := make(chan struct{})
	go func() {
		select {
		case <-sigChan:
			close(doneChan)
		case <-doneChan:
		}
	}()

	var group sync.WaitGroup
	group.Add(1)
	go func() {
		app.runWebServer(doneChan)
		group.Done()
	}()
	if app.levelDB != nil {
		group.Add(1)
		go func() {
			<-doneChan
			app.levelDB.Close()
			group.Done()
		}()
	}
	group.Wait()
	log.Println("Server shut down")
}
