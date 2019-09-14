// Redisscan continually scans the whole Redis keyspace and performs actions on it.
package redisscan

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/gomodule/redigo/redis"
	lredis "github.com/syzoj/syzoj-ng-go/lib/redis"
)

var log = logrus.StandardLogger()

// The interface for a redis scan handler.
type Redisscanner interface {
	RedisScan(context.Context, string)
}

// Represents a Redisscan instance.
type Redisscan struct {
	Redis *lredis.PoolWrapper
	Ratio float64 // how much CPU time to use, must be between 0 and 1
	BatchSize int
	Match string
	Handler Redisscanner

	cursor int64
}

// Creates a default redisscan instance with provided arguments.
func DefaultRedisscan(r *lredis.PoolWrapper, handle Redisscanner) *Redisscan {
	return &Redisscan{
		Redis: r,
		Ratio: 0.01,
		Handler: handle,
		BatchSize: 10,
	}
}

// Runs the redisscan instance.
func (r *Redisscan) Run(ctx context.Context) error {
	if r.Ratio < 0 || r.Ratio > 1 {
		panic("Redisscan: invalid ratio")
	}
	r.cursor = 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		startTime := time.Now()
		r.loop(ctx)
		dur := time.Now().Sub(startTime)
		time.Sleep(time.Duration(float64(dur) / r.Ratio))
	}
}

func (r *Redisscan) loop(ctx context.Context) {
	var (
		res []interface{}
		err error
	)
	if r.Match == "" {
		res, err = redis.Values(r.Redis.DoContext(ctx, "SCAN", r.cursor, "COUNT", r.BatchSize))
	} else {
		res, err = redis.Values(r.Redis.DoContext(ctx, "SCAN", r.cursor, "MATCH", r.Match, "COUNT", r.BatchSize))
	}
	if err != nil {
		log.Error(err)
		return
	}
	r.cursor, err = redis.Int64(res[0], nil)
	if err != nil {
		log.Error(err)
		return
	}
	keys, err := redis.Strings(res[1], nil)
	if err != nil {
		log.Error(err)
		return
	}
	for _, key := range keys {
		r.Handler.RedisScan(ctx, key)
	}
}
