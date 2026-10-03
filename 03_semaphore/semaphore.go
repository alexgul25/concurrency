package semaphore

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Semaphore struct {
	permits  uint32
	waitings uint32
}

func New(n int) *Semaphore {
	if n < 0 {
		panic("semaphore permits can't be negative!!!")
	}
	return &Semaphore{permits: uint32(n)}
}

func (s *Semaphore) Acquire() {
	for {
		curPermits := atomic.LoadUint32(&s.permits)
		if curPermits != 0 && atomic.CompareAndSwapUint32(&s.permits, curPermits, curPermits-1) {
			return
		}
		if curPermits == 0 {
			atomic.AddUint32(&s.waitings, 1)
			futex.Wait(&s.permits, 0)
			atomic.AddUint32(&s.waitings, ^uint32(0))
		}
	}
}

func (s *Semaphore) TryAcquire() bool {
	for {
		curPermits := atomic.LoadUint32(&s.permits)
		if curPermits == 0 {
			return false
		}
		if atomic.CompareAndSwapUint32(&s.permits, curPermits, curPermits-1) {
			return true
		}
	}
}

func (s *Semaphore) Release() {
	atomic.AddUint32(&s.permits, 1)
	if atomic.LoadUint32(&s.waitings) != 0 {
		futex.Wake(&s.permits)
	}
}

func (s *Semaphore) Available() int {
	return int(atomic.LoadUint32(&s.permits))
}
