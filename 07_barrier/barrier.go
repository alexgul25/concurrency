package barrier

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Barrier struct {
	need    uint32
	arrived uint32
	round   uint32
}

func New(n int) *Barrier {
	if n <= 0 {
		panic("barrier size must be positive")
	}
	return &Barrier{need: uint32(n)}
}

func (b *Barrier) Wait() {
	if b.need == 0 {
		panic("barrier size must be positive")
	}

	curRound := atomic.LoadUint32(&b.round)

	newArrived := atomic.AddUint32(&b.arrived, 1)
	if newArrived == b.need {
		atomic.StoreUint32(&b.arrived, 0)
		atomic.AddUint32(&b.round, 1)
		futex.WakeAll(&b.round)
	}

	for curRound == atomic.LoadUint32(&b.round) {
		futex.Wait(&b.round, curRound)
	}
}
