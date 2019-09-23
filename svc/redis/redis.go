// redis wraps around github.com/gomodule/redigo/redis to provide convenience functions.
// lib/redis contains utility functions for easier interaction with github.com/gomodule/redigo/redis
package redis

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

// Redis service.
type RedisService struct {
	*redis.Pool

	// For redis keyspace scanning
	BatchSize    int
	Ratio        float64
	ErrorHandler func(error)

	root   node
	cursor int64
}

// Creates a Redis service with default settings.
func DefaultRedisService(p *redis.Pool) *RedisService {
	return &RedisService{
		Pool:      p,
		BatchSize: 10,
		Ratio:     0.01,
		ErrorHandler: func(err error) {
			log.Error("error while scanning Redis: %s", err)
		},
	}
}

// Do takes a connection from the pool and does the command.
func (s *RedisService) Do(cmdName string, args ...interface{}) (interface{}, error) {
	conn := s.Get()
	defer conn.Close()
	return conn.Do(cmdName, args...)
}

// DoContext takes a connection from the pool and does the command.
// Context deadline is respected.
func (s *RedisService) DoContext(ctx context.Context, cmdName string, args ...interface{}) (interface{}, error) {
	conn, err := s.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conni := conn.(redis.ConnWithTimeout)
	deadline, ok := ctx.Deadline()
	var timeout time.Duration
	if ok {
		timeout = deadline.Sub(time.Now())
		if timeout < 0 {
			return nil, ctx.Err()
		}
	}
	if timeout > 0 {
		return conni.DoWithTimeout(timeout, cmdName, args...)
	} else {
		return conni.Do(cmdName, args...)
	}
}

// WithCache queries redis for cached content before calling func.
func (s *RedisService) WithCache(ctx context.Context, key string, f func() ([]byte, time.Duration, error)) ([]byte, error) {
	val, err := redis.Bytes(s.DoContext(ctx, "GET", key))
	if err == redis.ErrNil {
		var d time.Duration
		val, d, err = f()
		if err != nil {
			return nil, err
		}
		_, err2 := s.DoContext(ctx, "SET", key, val, "PX", int64(d/time.Millisecond))
		if err2 != nil { // this is not critical error
			log.WithError(err2).Warning("failed to write redis cache")
		}
	}
	return val, err
}
