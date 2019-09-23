package redis

import (
	"context"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Redis timer is designed around a ZSET. The score is used as Unix time in milliseconds.
// It polls a Redis ZSET to fetch latest messages. Instead of popping items out of ZSET,
// it modifies the score to some time in the future so items can't be lost. The item
// will be deleted automatically after the associated work is done.
type RedisTimer struct {
	r         *RedisService
	Key       string
	Handler   func(context.Context, string) error
	Timeout   time.Duration
	Interval  time.Duration // Maximum interval between polls
	BatchSize int

	tmr     *time.Timer
	tmrmu   sync.Mutex
	nextRun time.Time
}

func (r *RedisService) DefaultTimer(key string, handler func(context.Context, string) error) *RedisTimer {
	return &RedisTimer{
		r:         r,
		Key:       key,
		Handler:   handler,
		Timeout:   time.Second * 60,
		Interval:  time.Second * 5,
		BatchSize: 10,
		tmr:       time.NewTimer(0),
	}
}

func TimeToMillisecond(t time.Time) int64 {
	return t.Unix()*1000 + int64(t.Nanosecond())/1000000
}

func TimeFromMillisecond(t int64) time.Time {
	return time.Unix(t/1000, (t%1000)*1000000)
}

// Fetch a few elements, set their score to a new value, then find the lowest score.
var timerFetchScript = NewScript([]byte(`
local keys = redis.call("ZRANGEBYSCORE", KEYS[1], "-inf", ARGV[1])
local args = {}
for i,key in pairs(keys) do
  args[#args+1] = ARGV[2]
  args[#args+1] = key
end
if next(args) ~= nil then redis.call("ZADD", KEYS[1], unpack(args)) end
local first = redis.call("ZRANGE", KEYS[1], 0, 0, "WITHSCORES")
if next(first) == nil then first = false else first = first[2] end
return { first, keys }
`))

func (t *RedisTimer) Run(ctx context.Context) error {
	t.reschedule(time.Now())
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.tmr.C:
		}
		now := time.Now()
		nex := now.Add(t.Timeout)
		res, err := redis.Values(t.r.EvalContext(ctx, timerFetchScript, []string{t.Key}, []interface{}{TimeToMillisecond(now), TimeToMillisecond(nex)}))
		if err != nil {
			log.WithField("key", t.Key).WithError(err).Errorf("failed to fetch redis timer")
			t.reschedule(time.Now().Add(t.Interval))
			continue
		}
		keys, err := redis.Strings(res[1], nil)
		if err != nil {
			panic(err)
		}
		for _, key := range keys {
			if err := t.Handler(ctx, key); err != nil {
				log.WithField("key", t.Key).WithError(err).Error("failed to handle redis timer")
			} else if err = t.Delete(ctx, key); err != nil {
				log.WithField("key", t.Key).WithError(err).Error("failed to drain redis timer")
			}
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
		nexTime := time.Now().Add(t.Interval)
		firstTime, err := redis.Int64(res[0], nil)
		if err == nil {
			ft := TimeFromMillisecond(firstTime)
			if ft.Before(nexTime) {
				nexTime = ft
			}
		} else if err != redis.ErrNil {
			panic(err)
		}
		t.reschedule(nexTime)
	}
}

func (t *RedisTimer) reschedule(v time.Time) {
	t.tmrmu.Lock()
	defer t.tmrmu.Unlock()
	d := v.Sub(time.Now())
	if d < 0 {
		d = 0
	}
	t.tmr.Reset(d)
	t.nextRun = v
}

func (t *RedisTimer) schedule(v time.Time) {
	t.tmrmu.Lock()
	defer t.tmrmu.Unlock()
	if v.After(t.nextRun) {
		return
	}
	d := v.Sub(time.Now())
	if d < 0 {
		d = 0
	}
	t.tmr.Reset(d)
	t.nextRun = v
}

// Schedule key at the specified time. If key already exists, it is rescheduled.
func (t *RedisTimer) Schedule(ctx context.Context, key string, at time.Time) error {
	t.schedule(at)
	_, err := t.r.DoContext(ctx, "ZADD", t.Key, TimeToMillisecond(at), key)
	return err
}

// Delete key.
func (t *RedisTimer) Delete(ctx context.Context, key string) error {
	_, err := t.r.DoContext(ctx, "ZREM", t.Key, key)
	return err
}
