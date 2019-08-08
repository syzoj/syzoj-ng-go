package config

import (
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

func OpenRedis(redisName string) (*redis.Pool, error) {
	addr := os.Getenv(redisName + "_REDIS_ADDR")
	pool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
	return pool, nil
}
