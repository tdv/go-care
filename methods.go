// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

import (
	"sync"
	"time"
)

// 'Methods' represents an interface that allows to define
// the service's methods which responses you want to cache.
type Methods interface {
	// method - full method name
	//
	// Returns true and caching timeout if the method is found,
	// otherwise false and timeout in this case does not matter.
	Cacheable(method string) (bool, time.Duration)

	// Allows to add method for caching.
	// Returns the 'Methods' in order to be
	// more convenient methods adding like a chain.
	Add(method string, ttl time.Duration) Methods
	// Remove the method from allowed to be cached.
	Remove(method string)
	// Removes all methods.
	Clean()
}

type methodsStorage struct {
	mtx     sync.RWMutex
	methods map[string]time.Duration
}

func (s *methodsStorage) Cacheable(method string) (bool, time.Duration) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	ttl, ok := s.methods[method]
	return ok, ttl
}

func (s *methodsStorage) Add(method string, ttl time.Duration) Methods {
	s.mtx.Lock()
	s.mtx.Unlock()
	s.methods[method] = ttl
	return s
}

func (s *methodsStorage) Remove(method string) {
	s.mtx.Lock()
	s.mtx.Unlock()
	delete(s.methods, method)
}

func (s *methodsStorage) Clean() {
	newMap := make(map[string]time.Duration)
	s.mtx.Lock()
	s.mtx.Unlock()
	s.methods = newMap
}

func newMethodsStorage() Methods {
	return &methodsStorage{
		methods: make(map[string]time.Duration),
	}
}
