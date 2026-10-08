package rwmutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

// Биты поля RWMutex.state. Младшие 30 бит хранят число активных
// читателей, два старших бита — флаги писателя. Флаг ожидания
// не даёт новым читателям входить, чтобы писатель не голодал.
const (
	// writer выставлен, пока замок держит писатель.
	writer = 1 << 31

	// writerWaiting выставлен, пока писатель ждёт освобождения замка.
	writerWaiting = 1 << 30
)

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
		panic("RUnlock without RLock!!!")
	}

	if newState == writerWaiting {
		futex.WakeAll(&rw.state)
	}
}

func (rw *RWMutex) Lock() {
	for {
		curState := atomic.LoadUint32(&rw.state)

		if curState&^writerWaiting == 0 {
			if atomic.CompareAndSwapUint32(&rw.state, curState, writer) {
				return
			}
		} else if atomic.CompareAndSwapUint32(&rw.state, curState, curState|writerWaiting) {
			futex.Wait(&rw.state, curState|writerWaiting)
		}

	}
}

func (rw *RWMutex) Unlock() {
	for {
		curState := atomic.LoadUint32(&rw.state)
		if curState&writer == 0 {
			panic("Unlock without Lock!!!")
		}

		if atomic.CompareAndSwapUint32(&rw.state, curState, 0) {
			futex.WakeAll(&rw.state)
			return
		}
	}
}
