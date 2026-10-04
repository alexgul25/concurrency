package rwmutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const writer = 1 << 31
const writerWaiting = 1 << 30

type RWMutex struct {
	state uint32
}

func (rw *RWMutex) RLock() {
	for {
		curState := atomic.LoadUint32(&rw.state)
		if curState&writer == 0 && curState&writerWaiting == 0 && atomic.CompareAndSwapUint32(&rw.state, curState, curState+1) {
			return
		}
		if curState&writer != 0 || curState&writerWaiting != 0 {
			futex.Wait(&rw.state, curState)
		}
	}
}

func (rw *RWMutex) RUnlock() {
	newState := atomic.AddUint32(&rw.state, ^uint32(0))
	oldState := newState + 1

	if oldState&^(writer|writerWaiting) == 0 {
		panic("state has gone into the negative!!!")
	}

	if newState == 0 || newState == writerWaiting {
		futex.WakeAll(&rw.state)
	}
}

func (rw *RWMutex) Lock() {
	for {
		curState := atomic.LoadUint32(&rw.state)

		if curState == 0 && atomic.CompareAndSwapUint32(&rw.state, curState, writer) {
			return
		}

		if curState == writerWaiting && atomic.CompareAndSwapUint32(&rw.state, curState, writer) {
			return
		}

		curState = atomic.OrUint32(&rw.state, writerWaiting)

		if curState != writerWaiting {
			futex.Wait(&rw.state, curState)
		}
	}
}

func (rw *RWMutex) Unlock() {
	for {
		curState := atomic.LoadUint32(&rw.state)
		if curState&writer == 0 {
			panic("Unlock without lock!!!")
		}

		if atomic.CompareAndSwapUint32(&rw.state, curState, 0) {
			futex.WakeAll(&rw.state)
			return
		}
	}
}
