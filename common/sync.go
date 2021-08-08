package common

import (
	"sync"
)

// Sync struct.
type Sync struct {
	done chan struct{}
	err  error
	mu   sync.Mutex
}

// Lock locks mutex.
func (s *Sync) Lock() {
	s.mu.Lock()
}

// Unlock unlocks mutex.
func (s *Sync) Unlock() {
	s.mu.Unlock()
}

// Init initializes Sync instance.
func (s *Sync) Init() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.done != nil {
		return false
	}
	s.done = make(chan struct{})
	return true
}

// Done returns a channel that closes when the job is finished.
func (s *Sync) Done() <-chan struct{} {
	s.mu.Lock()
	done := s.done
	s.mu.Unlock()
	return done
}

// Wait blocks goroutine until completion.
func (s *Sync) Wait() {
	<-s.Done()
}

// Close ...
func (s *Sync) Close() {
	s.ErrClose(nil)
}

// LockClose ...
func (s *Sync) LockClose(f func() error) {
	s.mu.Lock()
	s.close(f())
	s.mu.Unlock()
}

// ErrClose ...
func (s *Sync) ErrClose(err error) {
	s.mu.Lock()
	s.close(err)
	s.mu.Unlock()
}

func (s *Sync) close(err error) {
	if s.done != nil && s.err == nil {
		select {
		case <-s.done:
		default:
			s.err = err
			close(s.done)
		}
	}
}

// Err ...
func (s *Sync) Err() error {
	s.mu.Lock()
	e := s.err
	s.mu.Unlock()
	return e
}
