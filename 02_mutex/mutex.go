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

const spinsBeforeWait = 100

type Mutex struct {
	state uint32
}

func (m *Mutex) Lock() {
	for range spinsBeforeWait {
		if atomic.LoadUint32(&m.state) == free && atomic.CompareAndSwapUint32(&m.state, free, held) {
			return
		}
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
		panic("Unknown state")
	}
}
