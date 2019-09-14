// lib/redis contains utility functions for easier interaction with github.com/gomodule/redigo/redis
package redis

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

// PoolWrapper wraps around a *redis.Pool to provide command execution.
type PoolWrapper struct {
	*redis.Pool
}

func WrapPool(p *redis.Pool) *PoolWrapper {
	return &PoolWrapper{Pool: p}
}

// Do takes a connection from the pool and does the command.
func (w *PoolWrapper) Do(cmdName string, args ...interface{}) (interface{}, error) {
	conn := w.Get()
	defer conn.Close()
	return conn.Do(cmdName, args...)
}

// DoContext takes a connection from the pool and does the command.
// Context deadline is respected.
func (w *PoolWrapper) DoContext(ctx context.Context, cmdName string, args ...interface{}) (interface{}, error) {
	conn, err := w.GetContext(ctx)
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
func WithCache(ctx context.Context, pool *PoolWrapper, key string, f func() ([]byte, error)) ([]byte, error) {
	val, err := redis.Bytes(pool.DoContext(ctx, "GET", key))
	if err == redis.ErrNil {
		val, err = f()
		if err != nil {
			return nil, err
		}
		_, err2 := pool.DoContext(ctx, "SET", key, val, "EX", 86400)
		if err2 != nil { // this is not critical error
			log.WithError(err2).Warning("failed to write redis cache")
		}
	}
	return val, err
}
