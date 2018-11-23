package api

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/syzoj/syzoj-ng-go/app/util"
)

type Session struct {
	SessionId  string
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
	if val, ok := sess_map["user-id"]; ok {
		if userId, err := util.UUIDFromBytes([]byte(val)); err != nil {
			sess.AuthUserId = userId
		} else {
			panic(fmt.Sprintf("Invalid user-id field for session %s, err %s", sessId, err.Error()))
		}
	}
	return
}

func (srv *ApiServer) SaveSession(r *http.Request, w http.ResponseWriter, sess *Session) {
	var setExpire = false
	if sess.SessionId == "" {
		var buf [32]byte
		_, err := rand.Read(buf[:])
		if err != nil {
			panic(err)
		}

		sessId := make([]byte, 64)
		hex.Encode(sessId, buf[:])
		sess.SessionId = string(sessId)
		http.SetCookie(w, &http.Cookie{
			Name:     "SYZOJSESSION",
			Value:    sess.SessionId,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})
		setExpire = true
	}

	key := fmt.Sprintf("sess:%s", sess.SessionId)
	if _, err := srv.redis.Del(key).Result(); err != nil {
		panic(err)
	}
	rmap := make(map[string]interface{})
	if sess.IsLoggedIn() {
		rmap["user-id"] = sess.AuthUserId
	}
	if _, err := srv.redis.HMSet(key, rmap).Result(); err != nil {
		panic(err)
	}

	if setExpire {
		if _, err := srv.redis.Expire(key, 24*time.Hour).Result(); err != nil {
			panic(err)
		}
	}
}

func (sess *Session) IsLoggedIn() bool {
	return sess.AuthUserId == (util.UUID{})
}
