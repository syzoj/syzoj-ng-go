// http is a utility library for dealing with HTTP requests.
package http

import (
	"context"
	"time"

	"math/rand"
)

// M is a shortcut for map[string]interface{}.
type M map[string]interface{}

// RetryPolicy describes a retry strategy.
type RetryPolicy interface {
	// Keeps calling a function until the operation succeeds or retry limit exceeds.
	// It should return nil if the operation succeeds at least once. Otherwise, the
	// last error should be returned.
	Execute(context.Context, func() error) error
}

// ExpRetry tries the operation with an exponential backoff.
type ExpRetry struct {
	MaxRetry int           // Defaults to 5
	Unit     time.Duration // Defaults to 500*time.Millisecond
	Cap      time.Duration // Defaults to 15*time.Second
	Exp      float64       // Defaults to 2
	Jitter   float64       // Randomly varies the backoff to distribute requests evenly. Defaults to 0.4, must be [0,1)
}

// The default retry mechanism.
var DefaultRetry = &ExpRetry{}

// Implements the RetryPolicy interface.
func (r *ExpRetry) Execute(ctx context.Context, f func() error) (err error) {
	cnt := r.MaxRetry
	if cnt < 0 {
		panic("ExpRetry: negative MaxRetry")
	}
	if cnt == 0 {
		cnt = 5
	}
	d := r.Unit
	if d == 0 {
		d = 500 * time.Millisecond
	}
	jitter := r.Jitter
	if jitter < 0 || jitter > 1 {
		panic("ExpRetry: invalid jitter")
	}
	if jitter == 0 {
		jitter = 0.4
	}
	exp := r.Exp
	if exp < 0 {
		panic("ExpRetry: invalid exp")
	}
	if exp == 0 {
		exp = 2
	}
	for ; cnt > 0; cnt-- {
		if err = f(); err == nil {
			return nil
		}
		select {
		case <-time.After(time.Duration(float64(d) * (1 - rand.Float64()*jitter))):
		case <-ctx.Done():
			return ctx.Err()
		}
		d = time.Duration(float64(d) * exp)
	}
	return
}
