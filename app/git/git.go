package git

import (
    "net/http"
    "net/http/cgi"
    "os/exec"
    "bytes"
    "path/filepath"
    "log"
    "strings"
    "github.com/gorilla/mux"
    "database/sql"
)

type GitServer struct {
    cgiScript string
    db *sql.DB
    pathPrefix string
}

func CreateGitServer(db *sql.DB, pathPrefix string) (*GitServer, error) {
    srv := &GitServer{
        db: db,
        pathPrefix: pathPrefix,
    }
    dir, err := exec.Command("git", "--exec-path").Output()
    if err != nil {
        return nil, err
    }

    srv.cgiScript = filepath.Join(string(bytes.TrimRight(dir, "\r\n")), "git-http-backend")
    return srv, nil
}

func (srv *GitServer) Handle(w http.ResponseWriter, r *http.Request) {
    env := []string{
        "PATH_TRANSLATED=" + srv.pathPrefix + strings.TrimPrefix(r.URL.Path, "/git/"),
        "GIT_HTTP_EXPORT_ALL=",
    }
    vars := mux.Vars(r)
    gitId := vars["git-id"]
    authRequired := strings.HasSuffix(r.URL.Path, "/git-receive-pack")
    query := r.URL.Query()
    if services, ok := query["service"]; ok {
        if len(services) > 0 && services[0] == "git-receive-pack" {
            authRequired = true
        }
    }
    if authRequired {
        fail_auth := func(msg string) {
            w.Header().Add("WWW-Authenticate", "basic; realm=SYZOJ")
            http.Error(w, msg, 401)
        }
        internal_err := func(args ...interface{}) {
            log.Println(args...)
            http.Error(w, "Internal server error", 500)
        }

        user, pass, ok := r.BasicAuth()
        if !ok {
            fail_auth("Authorization required")
            return
        }

        rows, err := srv.db.Query("SELECT git_password FROM users WHERE user_name=$1", user)
        if err != nil {
            internal_err("Error querying database:", err)
            return
        }
        if !rows.Next() {
            fail_auth("User with specified username does not exist")
            return
        }

        var git_password *string
        rows.Scan(&git_password)
        if git_password == nil {
            fail_auth("User does not have git password")
            return
        }
        
        if pass != *git_password {
            fail_auth("Password incorrect")
            return
        }

        env = append(env, "REMOTE_USER=" + user)
    }

    var stderr bytes.Buffer
    handler := &cgi.Handler{
        Path: srv.cgiScript,
        Env: env,
        Stderr: &stderr,
    }
    handler.ServeHTTP(w, r)

    if stderr.Len() > 0 {
        log.Printf("[git-backend] Error at git %s: %s", gitId, stderr.String())
    }
}