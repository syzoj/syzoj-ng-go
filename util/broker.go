package util

import (
	"runtime"
	"sync"
)

type Subscriber interface {
	Notify() // Must not block
}

type ChanSubscriber struct {
	C chan struct{}
}

func NewChanSubscriber() ChanSubscriber {
	return ChanSubscriber{C: make(chan struct{}, 1)}
}
func (s ChanSubscriber) Notify() {
	select {
	case s.C <- struct{}{}:
	default:
	}
}

type Broker struct {
	mu     sync.Mutex
	ch     map[Subscriber]struct{}
	closed bool
	sem    chan struct{}
}

func NewBroker() *Broker {
	b := new(Broker)
	b.ch = make(map[Subscriber]struct{})
	b.sem = make(chan struct{}, 1)
	go b.work()
	return b
}

// Marks the broker as a blackhole, future broadcasts will not be delivered
func (b *Broker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ch = nil
	b.closed = true
	select {
	case b.sem <- struct{}{}:
	default:
	}
}

func (b *Broker) Subscribe(ch Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.ch[ch] = struct{}{}
}

func (b *Broker) Unsubscribe(ch Subscriber) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	delete(b.ch, ch)
}

func (b *Broker) Broadcast() {
	select {
	case b.sem <- struct{}{}:
	default:
	}
}

// TODO: Better broadcast strategies
func (b *Broker) work() {
	for range b.sem {
		b.mu.Lock()
		if b.closed {
			return
		}
		for ch := range b.ch {
			ch.Notify()
		}
		b.mu.Unlock()
		runtime.Gosched()
	}
}
