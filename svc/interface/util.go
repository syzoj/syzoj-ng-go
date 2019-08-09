package main

import (
	"context"
	"math"
	"time"

	"github.com/gomodule/redigo/redis"
)

func (app *App) tryGetCache(ctx context.Context, redisKey string) ([]byte, error) {
	conn, err := app.redisCache.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return redis.Bytes(conn.Do("LINDEX", redisKey, 0))
}

func (app *App) blockingGetCache(ctx context.Context, redisKey string, timeout time.Duration) ([]byte, error) {
	conn, err := app.redisCache.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	seconds := int(math.Ceil(float64(timeout) / float64(time.Second)))
	return redis.Bytes(conn.Do("BRPOPLPUSH", redisKey, redisKey, seconds))
}

func (app *App) waitForCache(ctx context.Context, redisKey string, timeout time.Duration, trigger func()) ([]byte, error) {
	val, err := app.tryGetCache(ctx, redisKey)
	if err == nil {
		return val, nil
	} else if err != redis.ErrNil {
		return nil, err
	}
	trigger()
	return app.blockingGetCache(ctx, redisKey, timeout)
}
