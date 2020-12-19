package multiplexer

import "sync"

type safeBool struct {
	sync.RWMutex
	v bool
}

func (s *safeBool) get() bool {
	s.RLock()
	v := s.v
	s.RUnlock()
	return v
}

func (s *safeBool) set(v bool) {
	s.Lock()
	s.v = v
	s.Unlock()
}
