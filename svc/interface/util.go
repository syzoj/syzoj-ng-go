package main

import (
	"context"
	"math"
	"time"

	"github.com/gomodule/redigo/redis"
)

func (app *App) tryGetCache(ctx context.Context, redisKey string) (string, error) {
	conn, err := app.redisCache.GetContext(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return redis.String(conn.Do("LINDEX", redisKey, 0))
}

func (app *App) blockingGetCache(ctx context.Context, redisKey string, timeout time.Duration) (string, error) {
	conn, err := app.redisCache.GetContext(ctx)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	seconds := int(math.Ceil(float64(timeout) / float64(time.Second)))
	return redis.String(conn.Do("BRPOPLPUSH", redisKey, redisKey, seconds))
}

func (app *App) waitForCache(ctx context.Context, redisKey string, timeout time.Duration, trigger func()) (string, error) {
	val, err := app.tryGetCache(ctx, redisKey)
	if err == nil {
		return val, nil
	} else if err != redis.ErrNil {
		return "", err
	}
	trigger()
	return app.blockingGetCache(ctx, redisKey, timeout)
}
