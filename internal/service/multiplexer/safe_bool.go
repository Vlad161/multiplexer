package multiplexer

import (
	"sync"
	"sync/atomic"
)

// SafeBoolMutex
type SafeBoolMutex struct {
	sync.RWMutex
	v bool
}

func (s *SafeBoolMutex) Get() bool {
	s.RLock()
	v := s.v
	s.RUnlock()
	return v
}

func (s *SafeBoolMutex) Set(v bool) {
	s.Lock()
	s.v = v
	s.Unlock()
}

// SafeBoolAtomic
type SafeBoolAtomic struct {
	v uint32
}

func (s *SafeBoolAtomic) Get() bool {
	return atomic.LoadUint32(&s.v) == 1
}

func (s *SafeBoolAtomic) Set(v bool) {
	var newUint uint32 = 0
	var oldUint uint32 = 1
	if v {
		newUint = 1
		oldUint = 0
	}
	atomic.CompareAndSwapUint32(&s.v, oldUint, newUint)
}

// SafeBoolAtomicType
type safeBoolAtomicType struct {
	v atomic.Value
}

func NewSafeBoolAtomicType() *safeBoolAtomicType {
	sf := safeBoolAtomicType{}
	sf.v.Store(false)
	return &sf
}

func (s *safeBoolAtomicType) Get() bool {
	return s.v.Load().(bool)
}

func (s *safeBoolAtomicType) Set(v bool) {
	s.v.Store(v)
}
