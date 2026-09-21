package once

import (
	"primitives/internal/futex"
	"sync/atomic"
)

const (
	ready = iota
	carried
	completed
)

type Once struct {
	state uint32
}

func (o *Once) Do(f func()) {
	if atomic.CompareAndSwapUint32(&o.state, ready, carried) {
		func() {
			defer func() { recover() }()
			f()
		}()
		atomic.StoreUint32(&o.state, completed)
		futex.WakeAll(&o.state)
		return
	}

	futex.Wait(&o.state, carried)
}

func (o *Once) Done() bool {
	return atomic.LoadUint32(&o.state) == completed
}
