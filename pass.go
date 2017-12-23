package spass

import (
	"sync"
	"time"

	"github.com/tmthrgd/go-memset"
)

// Structure that will store
// the password in safe space (like feminists)
type Shadow struct {
	// password storage
	chars []byte
	// for concurrency access
	locker sync.RWMutex
}

// Reads password and stores in secured space
func (s *Shadow) Read(cleanTime time.Duration) (err error) {
	s.locker.Lock()
	defer s.locker.Unlock()

	var chars []byte

	chars, err = read()
	if err == nil {
		defer s.CleanBytes(chars)
		// After minute, password will
		// be deleted if noclean was not setted

		copy(s.chars, chars)
		if cleanTime > 0 {
			go func() {
				<-time.After(cleanTime)
				s.Clean()
			}()
		}
	}
	return
}

// Get returns a pointer to password
func (s *Shadow) Get() *[]byte {
	s.locker.Lock()
	defer s.locker.Unlock()
	return &s.chars
}

// Clean deletes password replacing
// characters with 0's
func (s *Shadow) Clean() {
	s.locker.Lock()
	defer s.locker.Unlock()

	memset.Memset(s.chars, 0)
}

// Store saves parsed password inside
// Shadow structure
func (s *Shadow) Store(pass []byte) {
	s.locker.Lock()
	defer s.locker.Unlock()

	memset.Memset(s.chars, 0)
	s.chars = append(s.chars[:0], pass...)
}

func (s *Shadow) StoreString(pass string) {
	s.Store([]byte(pass))
}

func (s *Shadow) CleanBytes(chs []byte) {
	memset.Memset(chs, 0)
}
