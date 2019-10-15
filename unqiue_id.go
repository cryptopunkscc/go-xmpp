package xmpp

import (
	"strconv"
	"sync"
)

// UniqueID holds the state of a unique id generator
type UniqueID struct {
	lastID int
	mu     sync.Mutex
}

// Next returns the next unique ID
func (s *UniqueID) Next() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastID++
	return strconv.Itoa(s.lastID)
}
