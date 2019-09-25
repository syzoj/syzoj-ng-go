package redis

import (
	"context"
	"sync"

	"github.com/gomodule/redigo/redis"
)

// Pipeline wraps a Redis connection to provide pipelining utilities.
// It supports multiple concurrent callers.
type Pipeline struct {
	conn redis.Conn
	mu   sync.Mutex
	ch   chan interface{}
}

// Create a new pipeline. The pipeline must be closed after use.
func (r *RedisService) NewPipeline(ctx context.Context) (*Pipeline, error) {
	conn, err := r.Pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	p := &Pipeline{conn: conn, ch: make(chan interface{}, 1000)}
	go p.run()
	return p, nil
}

// Add a command to pipeline. Note that callback will not receive connection errors, only redis errors.
// Callback can be nil, in which case it will be ignored.
// The callback will be called in order in a dedicated goroutine.
func (p *Pipeline) Do(callback func(interface{}, error), cmdName string, args ...interface{}) {
	p.mu.Lock()
	p.conn.Send(cmdName, args...)
	p.ch <- callback
	p.mu.Unlock()
}

// Flush the pipeline and wait until all previous commands have returned.
func (p *Pipeline) Flush(ctx context.Context) error {
	p.mu.Lock()
	if err := p.conn.Flush(); err != nil {
		p.mu.Unlock()
		return err
	}
	p.mu.Unlock()
	ch := make(chan error)
	p.ch <- ch
	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		return err
	}
}

func (p *Pipeline) run() {
	var perr error
	for f := range p.ch {
		switch obj := f.(type) {
		case chan error: // Flush command
			obj <- perr
		case func(interface{}, error):
			data, err := p.conn.Receive()
			if err != nil {
				if _, ok := err.(redis.Error); !ok {
					perr = err
				}
			}
			if perr == nil && obj != nil {
				obj(data, err)
			}
		}
	}
}

func (p *Pipeline) Close() error {
	close(p.ch)
	return p.conn.Close()
}

// A wrapper class to provide a callback that collects the results.
type RedisResult struct {
	Result interface{}
	Err    error
}

func (r *RedisResult) Save(a interface{}, e error) {
	r.Result = a
	r.Err = e
}
