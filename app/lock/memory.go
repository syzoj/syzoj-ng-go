package lock

import (
	"context"
	"sync"
	"time"
)

type memoryLockManager struct {
	data sync.Map
}

func CreateMemoryLockManager() LockManager {
	return new(memoryLockManager)
}

func (m *memoryLockManager) WithLockShared(ctx context.Context, id string, wait bool, handler func(context.Context, SharedLock) error) error {
	val, ok := m.data.Load(id)
	if !ok {
		val, _ = m.data.LoadOrStore(id, new(sync.RWMutex))
	}
	lock := val.(*sync.RWMutex)
	lock.RLock()
	defer lock.RUnlock()

	ctx1, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	return handler(ctx1, nil)
}

func (m *memoryLockManager) WithLockExclusive(ctx context.Context, id string, wait bool, handler func(context.Context, ExclusiveLock) error) error {
	val, ok := m.data.Load(id)
	if !ok {
		val, _ = m.data.LoadOrStore(id, new(sync.RWMutex))
	}
	lock := val.(*sync.RWMutex)
	lock.Lock()
	defer lock.Unlock()

	ctx1, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	return handler(ctx1, nil)
}
