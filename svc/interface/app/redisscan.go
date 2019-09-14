package app

import (
	"strings"
	"context"
)

// Implements the Redisscanner interface.
func (a *App) RedisScan(ctx context.Context, key string) {
	ind := strings.Index(key, ":")
	if ind == -1 {
		return
	}
	prefix := key[:ind]
	switch prefix {
	// case "session":
	case "stats":
		a.stats.RedisScan(ctx, key)
	}
}
