package redis

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

type XREAD struct {
	Key  string
	Msgs []*XREADMsg
}

type XREADMsg struct {
	ID   string
	Data map[string]string
}

// Parse XREAD and XREADGROUP result.
func ParseXREAD(data interface{}, err error) ([]*XREAD, error) {
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, redis.ErrNil
	}
	vals := data.([]interface{})
	res := make([]*XREAD, len(vals))
	for i, val := range vals {
		valv := val.([]interface{})
		xread := &XREAD{}
		if xread.Key, err = redis.String(valv[0], nil); err != nil {
			return nil, err
		}
		msgs := valv[1].([]interface{})
		xread.Msgs = make([]*XREADMsg, len(msgs))
		for j, msg := range msgs {
			msgv := msg.([]interface{})
			xmsg := &XREADMsg{}
			if xmsg.ID, err = redis.String(msgv[0], nil); err != nil {
				return nil, err
			}
			kvs, err := redis.Values(msgv[1], nil)
			if err != nil {
				return nil, err
			}
			xmsg.Data = make(map[string]string, len(kvs)/2)
			for i := 0; i < len(kvs); i += 2 {
				k, err := redis.String(kvs[i], nil)
				if err != nil {
					return nil, err
				}
				v, err := redis.String(kvs[i+1], nil)
				if err != nil {
					return nil, err
				}
				xmsg.Data[k] = v
			}
			xread.Msgs[j] = xmsg
		}
		res[i] = xread
	}
	return res, nil
}

// Read up to <count> entries from a Redis 5 stream until context is done.
// If count==0, it reads as many entries as possible.
func (r *RedisService) ReadStreamOnce(ctx context.Context, key string, lastId string, count int) ([]*XREADMsg, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		var (
			data []*XREAD
			err  error
		)
		if count == 0 {
			data, err = ParseXREAD(r.DoContext(ctx, "XREAD", "BLOCK", 60000, "STREAMS", key, lastId))
		} else {
			data, err = ParseXREAD(r.DoContext(ctx, "XREAD", "COUNT", count, "BLOCK", 60000, "STREAMS", key, lastId))
		}
		if err == redis.ErrNil {
			continue
		}
		return data[0].Msgs, err
	}
}

// Continuously reads from a single Redis 5 stream until context is done.
// When the operation completes, an error is returned, and both channels will be closed.
// Be sure to read from the returned channels or the stream will block.
// Set lastId to "0" to read from beginning, or "$" to read new entries (racy).
func (r *RedisService) ReadStream(ctx context.Context, key string, lastId string) (<-chan *XREADMsg, <-chan error) {
	msgCh := make(chan *XREADMsg)
	msgErr := make(chan error, 1)
	go func() {
		var err error
		defer func() {
			msgErr <- err
			close(msgCh)
			close(msgErr)
		}()
		for {
			var data []*XREADMsg
			data, err = r.ReadStreamOnce(ctx, key, lastId, 0)
			if err != nil {
				return
			}
			for _, msg := range data {
				lastId = msg.ID
				select {
				case <-ctx.Done():
					err = ctx.Err()
					return
				case msgCh <- msg:
				}
			}
		}
	}()
	return msgCh, msgErr
}

// Reads from a consumer group in a stream, using sema as the semaphore. It receives struct{}{} as "tokens" and reads a message for every token.
func (r *RedisService) ReadStreamGroup(ctx context.Context, key string, group string, consumer string, sema <-chan struct{}) (<-chan *XREADMsg, <-chan error) {
	msgCh := make(chan *XREADMsg)
	errCh := make(chan error, 1)
	go func() {
		var count int
		var err error
		lastId := "0"
		defer func() {
			errCh <- err
			close(msgCh)
			close(errCh)
		}()
		for {
			if count == 0 {
				select {
				case <-ctx.Done():
					return
				case <-sema:
					count++
				}
			}
			var data []*XREAD
			data, err = ParseXREAD(r.DoContext(ctx, "XREADGROUP", "GROUP", group, consumer, "COUNT", 1, "BLOCK", 60000, "STREAMS", key, lastId))
			if err == redis.ErrNil {
				err = nil
				continue
			}
			if err != nil {
				return
			}
			if len(data) == 0 || len(data[0].Msgs) == 0 {
				lastId = ">"
				continue
			}
			for _, msg := range data[0].Msgs {
				if lastId != ">" {
					lastId = msg.ID
				}
				msgCh <- msg
				count--
			}
		}
	}()
	return msgCh, errCh
}
