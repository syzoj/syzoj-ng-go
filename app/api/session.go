package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/google/uuid"

	"github.com/syzoj/syzoj-ng-go/app/lock"
	model_session "github.com/syzoj/syzoj-ng-go/app/model/session"
)

var SessionNotFoundError = errors.New("Session not found")

func (srv *ApiServer) createSession(ctx context.Context) (sessId uuid.UUID, err error) {
	sessId = uuid.New()
	key := fmt.Sprintf("session:%s", sessId)
	data, err1 := json.Marshal(model_session.Session{})
	if err1 != nil {
		panic(err1)
	}

	err = srv.lockManager.WithLockExclusive(ctx, key, func(ctx context.Context, l lock.ExclusiveLock) error {
		_, err := srv.redis.SetNX(key, data, 0).Result()
		return err
	})
	return
}

func (srv *ApiServer) withSessionShared(ctx context.Context, sessId uuid.UUID, handler func(context.Context, *model_session.Session) error) error {
	key := fmt.Sprintf("session:%s", sessId)
	return srv.lockManager.WithLockShared(ctx, key, func(ctx context.Context, l lock.SharedLock) error {
		data, err := srv.redis.Get(key).Result()
		if err != nil {
			return err
		}
		sess := new(model_session.Session)
		if err := json.Unmarshal([]byte(data), sess); err != nil {
			return err
		}
		return handler(ctx, sess)
	})
}

func (srv *ApiServer) withSessionExclusive(ctx context.Context, sessId uuid.UUID, handler func(context.Context, *model_session.Session) error) error {
	key := fmt.Sprintf("session:%s", sessId)
	return srv.lockManager.WithLockExclusive(ctx, key, func(ctx context.Context, l lock.ExclusiveLock) error {
		data, err := srv.redis.Get(key).Result()
		if err != nil {
			return err
		}
		sess := new(model_session.Session)
		if err := json.Unmarshal([]byte(data), sess); err != nil {
			return err
		}
		return handler(ctx, sess)
	})
}

// Tries to get the current session and lock it for use.
// If the session is not valid it will try to create another.
func (srv *ApiServer) WithSessionShared(ctx context.Context, w http.ResponseWriter, r *http.Request, handler func(ctx context.Context, sessId uuid.UUID, sess *model_session.Session) error) error {
	var sessId uuid.UUID
	if sessCookie, err := r.Cookie("SYZOJSESSION"); err != nil {
		sessId, _ = uuid.Parse(sessCookie.Value)
	}
	var success bool = false
	err := srv.withSessionShared(ctx, sessId, func(ctx context.Context, sess *model_session.Session) error {
		success = true
		return handler(ctx, sessId, sess)
	})
	if success {
		return err
	} else if err == redis.Nil {
		sessId, err := srv.createSession(ctx)
		if err != nil {
			return err
		} else {
			return srv.withSessionShared(ctx, sessId, func(ctx context.Context, sess *model_session.Session) error {
				return handler(ctx, sessId, sess)
			})
		}
	} else {
		return err
	}
}

// Tries to get the current session and lock it for use.
// If the session is not valid it will try to create another.
func (srv *ApiServer) WithSessionExclusive(ctx context.Context, w http.ResponseWriter, r *http.Request, handler func(ctx context.Context, sessId uuid.UUID, sess *model_session.Session) error) error {
	var sessId uuid.UUID
	if sessCookie, err := r.Cookie("SYZOJSESSION"); err == nil {
		sessId, _ = uuid.Parse(sessCookie.Value)
	}
	var success bool = false
	err := srv.withSessionExclusive(ctx, sessId, func(ctx context.Context, sess *model_session.Session) error {
		success = true
		return handler(ctx, sessId, sess)
	})
	if success {
		return err
	} else if err == redis.Nil {
		sessId, err := srv.createSession(ctx)
		if err != nil {
			return err
		} else {
			return srv.withSessionExclusive(ctx, sessId, func(ctx context.Context, sess *model_session.Session) error {
				return handler(ctx, sessId, sess)
			})
		}
	} else {
		return err
	}
}

func (srv *ApiServer) SaveSession(ctx context.Context, sessId uuid.UUID, sess *model_session.Session) error {
	key := fmt.Sprintf("session:%s", sessId)
	data, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	return srv.redis.Set(key, data, 0).Err()
}
