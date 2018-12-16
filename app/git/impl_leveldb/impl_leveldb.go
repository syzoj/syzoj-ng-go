package impl_leveldb

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/cgi"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/syndtr/goleveldb/leveldb"
	go_git "gopkg.in/src-d/go-git.v4"

	"github.com/syzoj/syzoj-ng-go/app/git"
)

var log = logrus.StandardLogger()

type service struct {
	gitExecPath string
	cgiScript   string
	router      http.Handler
	db          *leveldb.DB
	pathPrefix  string
	handler     map[string]git.GitHookHandler
	lock        sync.Mutex
}

func NewLevelDBGitService(db *leveldb.DB, pathPrefix string) (s git.GitService, err error) {
	srv := &service{db: db}
	if srv.pathPrefix, err = filepath.Abs(pathPrefix); err != nil {
		return
	}
	var dir []byte
	dir, err = exec.Command("git", "--exec-path").Output()
	if err != nil {
		return
	}
	srv.gitExecPath = string(bytes.TrimRight(dir, "\r\n"))
	srv.cgiScript = filepath.Join(srv.gitExecPath, "git-http-backend")
	srv.handler = make(map[string]git.GitHookHandler)

	gitRouter := mux.NewRouter()
	gitRouter.HandleFunc("/{git-id}/HEAD", srv.Handle).Methods("GET")
	gitRouter.HandleFunc("/{git-id}/info/refs", srv.Handle).Methods("GET")
	gitRouter.HandleFunc("/{git-id}/objects/info/alternates", srv.Handle).Methods("GET")
	gitRouter.HandleFunc("/{git-id}/objects/info/http-alternates", srv.Handle).Methods("GET")
	gitRouter.HandleFunc("/{git-id}/objects/info/packs", srv.Handle).Methods("GET")
	gitRouter.HandleFunc("/{git-id}/objects/{object-id-prefix:[0-9a-f]{2}}/{object-id-suffix:[0-9a-f]{38}}", srv.Handle).Methods("GET")
	gitRouter.HandleFunc("/{git-id}/objects/pack/pack-{pack-id:[0-9a-f]{40}}.pack", srv.Handle).Methods("GET")
	gitRouter.HandleFunc("/{git-id}/objects/pack/pack-{pack-id:[0-9a-f]{40}}.idx", srv.Handle).Methods("GET")
	gitRouter.HandleFunc("/{git-id}/git-upload-pack", srv.Handle).Methods("POST")
	gitRouter.HandleFunc("/{git-id}/git-receive-pack", srv.Handle).Methods("POST")
	srv.router = http.StripPrefix("/git", gitRouter)
	s = srv
	return
}

func (srv *service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.router.ServeHTTP(w, r)
}

func (srv *service) Handle(w http.ResponseWriter, r *http.Request) {
	var err error

	var vars = mux.Vars(r)
	var _gitId = vars["git-id"]
	var gitId uuid.UUID
	if gitId, err = uuid.Parse(_gitId); err != nil {
		log.Warning("Git: failed to parse uuid")
		return
	}
	var env = []string{
		"GIT_PROJECT_ROOT=" + srv.pathPrefix,
		"GIT_HTTP_EXPORT_ALL=",
	}

	var authRequired bool
	if strings.HasSuffix(r.URL.Path, "/git-receive-pack") {
		authRequired = true
	}
	if services, ok := r.URL.Query()["service"]; ok {
		if len(services) > 0 && services[0] == "git-receive-pack" {
			authRequired = true
		}
	}

	var info *gitRepoInfo
	if info, err = srv.getRepoInfo(srv.db, gitId); err != nil {
		log.Warningf("Attempt to access nonexistent git repo: %s", gitId)
		return
	}

	if authRequired {
		w.Header().Add("WWW-Authenticate", "basic; realm=SYZOJ")
		user, _, ok := r.BasicAuth()
		if !ok || (info.Token == "" || user != info.Token) {
			http.Error(w, "Authentication required", 401)
			return
		}

		env = append(env, fmt.Sprintf("REMOTE_USER=%s", user))
	}

	var stderr bytes.Buffer
	handler := &cgi.Handler{
		Path:   srv.cgiScript,
		Env:    env,
		Stderr: &stderr,
	}
	log.Printf("data: %+v %+v\n", env, r)
	handler.ServeHTTP(w, r)

	if stderr.Len() > 0 {
		log.Warningf("[git-backend] Error at git cgi %s: %s", gitId, stderr.String())
	}
}

func (srv *service) AttachHookHandler(HookType string, Handler git.GitHookHandler) {
	if HookType == "" {
		panic("HookType cannot be empty")
	}
	if _, ok := srv.handler[HookType]; ok {
		panic("Attaching a hook handler twice")
	}
	srv.handler[HookType] = Handler
}

func (srv *service) CreateRepository(HookType string) (id uuid.UUID, err error) {
	if id, err = uuid.NewRandom(); err != nil {
		return
	}

	if _, err = go_git.PlainInit(filepath.Join(srv.pathPrefix, id.String()), true); err != nil {
		return
	}
	var info = gitRepoInfo{
		HookType: HookType,
	}
	if err = srv.putRepoInfo(srv.db, id, &info); err != nil {
		return
	}
	return
}

func (srv *service) ResetToken(id uuid.UUID) (token string, err error) {
	var bytes [8]byte
	if _, err = rand.Read(bytes[:]); err != nil {
		return
	}
	token = hex.EncodeToString(bytes[:])

	srv.lock.Lock()
	defer srv.lock.Unlock()
	var info *gitRepoInfo
	if info, err = srv.getRepoInfo(srv.db, id); err != nil {
		return
	}
	info.Token = token
	if err = srv.putRepoInfo(srv.db, id, info); err != nil {
		return
	}
	return
}

func (srv *service) Close() error {
	return nil
}
