package lock

import "context"

type LockManager interface {
	// Locks the specified resource in shared mode.
	// The inner context will be cancelled when the lock expires.
	// The handler may not be called at all, in which case an error will be returned.
	WithLockShared(ctx context.Context, id string, handler func(context.Context, SharedLock) error) error
	// Locks the specified resource in exclusive mode.
	// The inner context will be cancelled when the lock expires.
	// The handler may not be called at all, in which case an error will be returned.
	WithLockExclusive(ctx context.Context, id string, handler func(context.Context, ExclusiveLock) error) error
}

type SharedLock interface {
}

type ExclusiveLock interface {
}
