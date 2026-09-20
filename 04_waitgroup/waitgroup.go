package waitgroup

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type WaitGroup struct {
	count uint32
}

func (wg *WaitGroup) Add(delta int) {
	atomic.AddUint32(&wg.count, uint32(delta))
}

func (wg *WaitGroup) Done() {
	newCount := atomic.AddUint32(&wg.count, ^uint32(0))

	if newCount == ^uint32(0) {
		panic("count has overflowed or gone into the negative!!!")
	}

	if newCount == 0 {
		futex.WakeAll(&wg.count)
	}
}

func (wg *WaitGroup) Wait() {
	for {
		curCount := atomic.LoadUint32(&wg.count)
		if curCount == 0 {
			return
		}
		futex.Wait(&wg.count, curCount)
	}
}
