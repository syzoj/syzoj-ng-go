package config

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Creates a new *github.com/gomodule/redigo/redis.Pool instance from environment variables.
// The environment variable is ${prefix}REDIS_ADDR, in host:port format.
func NewRedis(prefix string) (*redis.Pool, error) {
	addr := os.Getenv(prefix + "REDIS_ADDR")
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
	return pool, nil
}
