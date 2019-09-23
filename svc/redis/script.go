package redis

import (
	"context"
	"crypto/sha1"
	"encoding/hex"

	"github.com/gomodule/redigo/redis"
)

// Script wraps a Lua script to be used in Redis.
type Script struct {
	Script []byte
	SHA1   [20]byte
	Hash   string
}

func NewScript(sc []byte) *Script {
	sum := sha1.Sum(sc)
	return &Script{
		Script: sc,
		SHA1:   sum,
		Hash:   hex.EncodeToString(sum[:]),
	}
}

// Evaluates a Lua script.
func (r *RedisService) EvalContext(ctx context.Context, sc *Script, keys []string, data []interface{}) (interface{}, error) {
	args := make([]interface{}, 2+len(keys)+len(data))
	args[0] = sc.Hash
	args[1] = len(keys)
	for i, key := range keys {
		args[i+2] = key
	}
	for i, dat := range data {
		args[i+2+len(keys)] = dat
	}
	val, err := r.DoContext(ctx, "EVALSHA", args...)
	if rerr, ok := err.(redis.Error); ok {
		if rerr[0:9] == "NOSCRIPT " {
			args[0] = sc.Script
			return r.DoContext(ctx, "EVAL", args...)
		}
	}
	return val, err
}
