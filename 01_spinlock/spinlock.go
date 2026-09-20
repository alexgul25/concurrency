package spinlock

import "sync/atomic"

type Spinlock struct {
	locked atomic.Bool
}

func (s *Spinlock) Lock() {
	for !s.locked.CompareAndSwap(false, true) {
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
	for s.locked.Load() {
	}

	for !s.locked.CompareAndSwap(false, true) {
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
