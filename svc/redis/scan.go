package redis

import (
	"context"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

type node struct {
	child map[string]*node
	scan  func(string)
}

// Handles a Redis prefix. The prefix must either be empty or end in a colon.
func (s *RedisService) ScanPrefix(prefix string, handler func(string)) {
	cur := &s.root
	for {
		ind := strings.Index(prefix, ":")
		if ind == -1 {
			break
		}
		name := prefix[:ind]
		prefix = prefix[ind+1:]
		if cur.scan != nil {
			panic("Scanner: conflicting prefix")
		}
		if cur.child == nil {
			cur.child = make(map[string]*node)
		}
		nex := &node{}
		cur.child[name] = nex
		cur = nex
	}
	if prefix != "" {
		panic("ScanPrefix: invalid prefix")
	}
	if cur.child != nil || cur.scan != nil {
		panic("ScanPrefix: conflicting prefix")
	}
	cur.scan = handler
}

func (r *RedisService) RunScanner(ctx context.Context) error {
	if r.Ratio < 0 || r.Ratio > 1 {
		panic("redis: invalid ratio")
	}
	r.cursor = 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		startTime := time.Now()
		if err := r.loop(ctx); err != nil {
			r.ErrorHandler(err)
		}
		dur := time.Now().Sub(startTime)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Duration(float64(dur) / r.Ratio)):
		}
	}
}

func (r *RedisService) loop(ctx context.Context) error {
	res, err := redis.Values(r.DoContext(ctx, "SCAN", r.cursor, "COUNT", r.BatchSize))
	if err != nil {
		return err
	}
	r.cursor, err = redis.Int64(res[0], nil)
	if err != nil {
		return err
	}
	keys, err := redis.Strings(res[1], nil)
	if err != nil {
		return err
	}
	for _, key := range keys {
		cur := &r.root
		for cur.child != nil {
			ind := strings.Index(key, ":")
			name := key[:ind]
			if nex, ok := cur.child[name]; ok {
				cur = nex
				key = key[ind+1:]
			} else {
				cur = nil
				break
			}
		}
		if cur != nil && cur.scan != nil {
			cur.scan(key)
		}
	}
	return nil
}
