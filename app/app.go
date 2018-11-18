package app

import (
    "net/http"
    "log"
    "os"
    "os/signal"
    "syscall"
	"database/sql"
    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
    "github.com/go-redis/redis"
	
	"github.com/syzoj/syzoj-ng-go/app/git"
	"github.com/syzoj/syzoj-ng-go/app/api"
)

type App struct {
    db *sql.DB
    redis *redis.Client
	httpServer *http.Server
	router *mux.Router
	gitServer *git.GitServer
	apiServer *api.ApiServer
}

func (app *App) SetupDB(conn string) error {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return err
	}
	app.db = db
	return nil
}

func (app *App) SetupRedis(options *redis.Options) error {
    app.redis = redis.NewClient(options)
    _, err := app.redis.Ping().Result()
    return err
}

func (app *App) SetupGitServer(gitPath string) error {
	server, err := git.CreateGitServer(app.db, gitPath)
	if err != nil {
		return err
	}

	app.gitServer = server
	return nil
}

func (app *App) SetupApiServer() error {
	server, err := api.CreateApiServer(app.db, app.redis)
	if err != nil {
		return err
	}

	app.apiServer = server
	return nil
}

func (app *App) SetupHttpServer(addr string) error {
	app.router = mux.NewRouter()
	app.httpServer = &http.Server{
		Addr: addr,
		Handler: app.router,
	}
	return nil
}

func (app *App) AddGitServer() {
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
}

func (app *App) AddApiServer() {
	apiRouter := app.router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/auth/register", app.apiServer.HandleAuthRegister).Methods("POST")
    apiRouter.HandleFunc("/auth/login", app.apiServer.HandleAuthLogin).Methods("POST")
    apiRouter.HandleFunc("/user/info", app.apiServer.HandleUserInfo).Methods("GET")
    apiRouter.HandleFunc("/group/create", app.apiServer.HandleGroupCreate).Methods("POST")
    apiRouter.PathPrefix("/").HandlerFunc(app.apiServer.HandleCatchAll)
}

func (app *App) Run() {
    errChan := make(chan error)
    go func() {
        log.Println("Starting web server at", app.httpServer.Addr)
        if err := app.httpServer.ListenAndServe(); err != nil {
            errChan <- err
        }
    }()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

    select {
        case err := <- errChan:
            log.Println("Web server error:", err)
        case sig := <- sigChan:
            log.Printf("Received signal %s, shutting down", sig)
            if err := app.httpServer.Shutdown(nil); err != nil {
                log.Println("Failed to shutdown server:", err)
            }
            log.Println("Web server shut down")
    }
}