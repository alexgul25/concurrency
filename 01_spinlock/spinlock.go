package spinlock

import (
	"runtime"
	"sync/atomic"
)

const spinsBeforeYield = 100

type Spinlock struct {
	locked atomic.Bool
}

func (s *Spinlock) Lock() {
	count := 0
	for !s.locked.CompareAndSwap(false, true) {
		count++
		if count == spinsBeforeYield {
			count = 0
			runtime.Gosched()
		}
	}
}

func (s *Spinlock) TryLock() bool {
	return s.locked.CompareAndSwap(false, true)
}

func (s *Spinlock) Unlock() {
	if s.locked.CompareAndSwap(true, false) {
		return
	}
	panic("Unlock without Lock!!!")
}

type TTAS struct {
	locked atomic.Bool
}

func (s *TTAS) Lock() {
	count := 0
	for {
		for s.locked.Load() {
			count++
			if count == spinsBeforeYield {
				count = 0
				runtime.Gosched()
			}
		}
		if s.locked.CompareAndSwap(false, true) {
			return
		}
	}
}

func (s *TTAS) TryLock() bool {
	return !s.locked.Load() && s.locked.CompareAndSwap(false, true)
}

func (s *TTAS) Unlock() {
	if s.locked.CompareAndSwap(true, false) {
		return
	}
	panic("Unlock without Lock!!!")
}
