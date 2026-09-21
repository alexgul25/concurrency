package rwmutex

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const writer = 1 << 31

type RWMutex struct {
	state uint32
}

/*

	Приведена наивная реализация

	Чтобы добиться справедливости, можно добавить в состояние, сигнализирующее о появлении писателя.
	Читатели обязаны проверять это состояние, и при необходимости пропускать писателя вперёд, тогда
	писателю достаточно дождаться только текущих чтений.

*/

func (rw *RWMutex) RLock() {
	for {
		curState := atomic.LoadUint32(&rw.state)
		if curState != writer && atomic.CompareAndSwapUint32(&rw.state, curState, curState+1) {
			return
		}
		futex.Wait(&rw.state, curState)
	}
}

func (rw *RWMutex) RUnlock() {
	newState := atomic.AddUint32(&rw.state, ^uint32(0))

	if newState == ^uint32(0) {
		panic("state has gone into the negative!!!")
	}

	if newState == 0 {
		futex.WakeAll(&rw.state)
	}
}

func (rw *RWMutex) Lock() {
	for {
		curState := atomic.LoadUint32(&rw.state)

		if curState == 0 && atomic.CompareAndSwapUint32(&rw.state, curState, writer) {
			return
		}

		futex.Wait(&rw.state, curState)
	}
}

func (rw *RWMutex) Unlock() {
	curState := atomic.LoadUint32(&rw.state)

	if curState != writer {
		panic("Unlock without lock!!!")
	}

	atomic.StoreUint32(&rw.state, 0)
	futex.WakeAll(&rw.state)
}
