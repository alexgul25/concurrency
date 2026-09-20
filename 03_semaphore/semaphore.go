package semaphore

import (
	"primitives/internal/futex"
	"sync/atomic"
)

type Semaphore struct {
	permits uint32
}

func New(n int) *Semaphore {
	return &Semaphore{permits: uint32(n)}
}

func (s *Semaphore) Acquire() {
	for {
		curPermits := atomic.LoadUint32(&s.permits)
		if curPermits != 0 && atomic.CompareAndSwapUint32(&s.permits, curPermits, curPermits-1) {
			return
		}
		futex.Wait(&s.permits, 0)
	}
}

func (s *Semaphore) TryAcquire() bool {
	curPermits := atomic.LoadUint32(&s.permits)
	return curPermits != 0 && atomic.CompareAndSwapUint32(&s.permits, curPermits, curPermits-1)
}

func (s *Semaphore) Release() {
	atomic.AddUint32(&s.permits, 1)
	futex.Wake(&s.permits)
}

func (s *Semaphore) Available() int {
	return int(atomic.LoadUint32(&s.permits))
}
