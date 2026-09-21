package mutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	free = iota
	held
	contended
)

type Mutex struct {
	state uint32
}

func (m *Mutex) Lock() {
	if atomic.CompareAndSwapUint32(&m.state, free, held) {
		return
	}

	oldState := atomic.SwapUint32(&m.state, contended)

	for oldState != free {
		futex.Wait(&m.state, contended)
		oldState = atomic.SwapUint32(&m.state, contended)
	}
}

func (m *Mutex) TryLock() bool {
	return atomic.CompareAndSwapUint32(&m.state, free, held)
}

func (m *Mutex) Unlock() {
	oldState := atomic.SwapUint32(&m.state, free)
	switch oldState {
	case free:
		panic("Unlock without lock!!!")
	case held:
		return
	case contended:
		futex.Wake(&m.state)
	default:
		panic("Unlnown state")
	}
}
