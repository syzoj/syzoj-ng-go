// The authentication middleware
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	lredis "github.com/syzoj/syzoj-ng-go/lib/redis"
)

// AuthMiddleware is a gin middleware for handling user auth.
type AuthMiddleware struct {
	Redis          *lredis.PoolWrapper // required
	CookieKey      string
	RedisKeyPrefix string
	GinKey         string
	Path           string
}

// A auth middleware with default settings.
func DefaultAuthMiddleware(redis *lredis.PoolWrapper) *AuthMiddleware {
	return &AuthMiddleware{
		Redis:          redis,
		CookieKey:      "SESSION",
		RedisKeyPrefix: "session:",
		GinKey:         "AUTH_INFO",
		Path:           "/",
	}
}

// AuthInfo contains information about currently authenticated user
type AuthInfo struct {
	UserId string
}

// The middleware function.
func (m *AuthMiddleware) Handle(c *gin.Context) {
	ctx := c.Request.Context()
	info := &AuthInfo{}
	sessKey, err := c.Cookie(m.CookieKey)
	if err == nil {
		reply, err := redis.String(m.Redis.DoContext(ctx, "GET", m.RedisKeyPrefix+sessKey))
		if err != nil {
			if err != redis.ErrNil {
				c.Error(err)
				reply = ""
			}
		}
		info.UserId = reply
	}
	c.Set(m.GinKey, info)
}

// Gets auth info from gin context.
func (m *AuthMiddleware) GetInfo(c *gin.Context) *AuthInfo {
	info, ok := c.Get(m.GinKey)
	if !ok {
		panic("auth: key doesn't exist")
	}
	return info.(*AuthInfo)
}

// Sets auth info.
func (m *AuthMiddleware) SetInfo(c *gin.Context, userId string, expire time.Duration) error {
	ctx := c.Request.Context()
	sessKey := makeRandomKey()
	_, err := m.Redis.DoContext(ctx, "SET", m.RedisKeyPrefix+sessKey, userId)
	if err != nil {
		return err
	}
	c.SetCookie("SESSION", sessKey, int(expire/time.Second)+1, "/", "", false, true)
	return nil
}

func makeRandomKey() string {
	var b [16]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
