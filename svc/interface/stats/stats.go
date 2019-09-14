// stats deals with statistical data.
// Such data can be reduced easily, does not need strong accuracy but may need to handle lots of data.
// The package does not emit errors. Instead, when it encounters an error, it generates a warning to the logger.
// The stats can be expected to be accurate if no warnings are generated.
package stats

import (
	"context"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	lredis "github.com/syzoj/syzoj-ng-go/lib/redis"
)

var log = logrus.StandardLogger()

// The stats service.
type Stats struct {
	Redis     *lredis.PoolWrapper
	KeyPrefix string

	// UpstreamCounter receives aggregated sums and save them.
	UpstreamCounter func(context.Context, string, int64)
}

// Increases a counter.
func (s *Stats) Inc(ctx context.Context, key string, num int64) {
	go func() {
		_, err := s.Redis.DoContext(ctx, "INCRBY", s.KeyPrefix+"cnt:"+key, num)
		if err != nil {
			log.WithError(err).Warning("failed to update stats")
		}
	}()
}

// Flush the current value upstream.
func (s *Stats) FlushInc(ctx context.Context, key string) {
	conn, err := s.Redis.GetContext(ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	rkey := s.KeyPrefix + "cnt:" + key
	conn.Send("MULTI")
	conn.Send("GET", rkey)
	conn.Send("DEL", rkey)
	r, err := conn.Do("EXEC")
	if err != nil {
		log.WithError(err).Warning("failed to reset counter")
		return
	}

	arr, err := redis.Values(r, nil)
	if err != nil {
		log.WithError(err).Warning("failed to parse return value")
		return
	}
	val, err := redis.Int64(arr[0], nil)
	if err != nil {
		log.WithError(err).Warning("failed to parse return value")
		return
	}
	s.UpstreamCounter(ctx, key, val)
}

// Implements redisscan.Redisscanner interface. Ignores keys that does not have the given prefix.
func (s *Stats) RedisScan(ctx context.Context, key string) {
	if !strings.HasPrefix(key, s.KeyPrefix) {
		return
	}
	key = key[len(s.KeyPrefix):]
	ind := strings.Index(key, ":")
	if ind == -1 {
		return
	}
	switch key[:ind] {
	case "cnt":
		s.FlushInc(ctx, key[ind+1:])
	}
}
