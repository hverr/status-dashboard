package broadcaster

import "sync"

// Broadcaster emits signals to its listeners.
//
// Signals cannot include information.
type Broadcaster struct {
	listenersLock sync.RWMutex
	listeners     []chan bool
}

// New creates a new broadcaster with no listeners.
func New() *Broadcaster {
	return &Broadcaster{}
}

// Listen adds a listener to the broadcaster.
//
// The returned channel will be activated for the next signal. After one
// signal the channel will be closed.
//
// If the broadcaster is stopped before a next signal was emitted, the
// channel will be closed.
func (b *Broadcaster) Listen() chan bool {
	c := make(chan bool, 1)

	b.listenersLock.Lock()
	b.listeners = append(b.listeners, c)
	b.listenersLock.Unlock()

	return c
}

// Emit a signal and notify all listeners.
//
// After emitting an event, all listeners will be removed.
func (b *Broadcaster) Emit() {
	b.listenersLock.Lock()
	for _, c := range b.listeners {
		c <- true
	}
	b.listeners = nil
	b.listenersLock.Unlock()
}
