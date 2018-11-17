package api

import (
	"time"
	"encoding/hex"
	"crypto/rand"
	"log"
	"fmt"
	"net/http"
	"github.com/go-redis/redis"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type Session struct {
	SessionId string
	LoggedIn bool
	AuthUserId util.UUID
}

func (srv *ApiServer) GetSession(req *http.Request) (sess *Session) {
	sess = &Session{}
	cookie, err := req.Cookie("SYZOJSESSION")
	if cookie == nil {
		return
	}
	sessId := cookie.Value

	sess_map, err := srv.redis.HGetAll(fmt.Sprintf("sess:%s", sessId)).Result()
	if err == redis.Nil || len(sess_map) == 0 {
		return
	}

	sess.SessionId = sessId
	val, ok := sess_map["user-id"]
	if ok {
		userId, err := util.UUIDFromBytes([]byte(val))
		if err == nil {
			sess.LoggedIn = true
			sess.AuthUserId = userId
		} else {
			log.Printf("Warning: invalid user-id field for session %s, err %s\n", sessId, err.Error())
		}
	} else {
		sess.LoggedIn = false
	}
	return
}

func (srv *ApiServer) SaveSession(r *http.Request, w http.ResponseWriter, sess *Session) error {
	if sess.SessionId == "" {
		var buf [32]byte
		_, err := rand.Read(buf[:])
		if err != nil {
			return err
		}

		sessId := make([]byte, 64)
		hex.Encode(sessId, buf[:])
		sess.SessionId = string(sessId)
		http.SetCookie(w, &http.Cookie{
			Name: "SYZOJSESSION",
			Value: sess.SessionId,
			Expires: time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})
	}

	key := fmt.Sprintf("sess:%s", sess.SessionId)
	if _, err := srv.redis.Del(key).Result(); err != nil {
		return err
	}
	rmap := make(map[string]interface{})
	if sess.LoggedIn {
		rmap["user-id"] = sess.AuthUserId
	}
	log.Println(key, rmap)
	if _, err := srv.redis.HMSet(key, rmap).Result(); err != nil {
		return err
	}

	return nil
}