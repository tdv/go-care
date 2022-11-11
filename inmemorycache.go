// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

import (
	"errors"
	"sync"
	"time"
)

type node struct {
	key  string
	prev *node
	next *node
}

type value struct {
	item *node
	val  []byte
}

// This built-in implementation of the in-memory cache
// doesn't support eviction by the TTL. It has developed only
// for demo and small MVP. In production, you need to use
// go-care with Redis, Memcached, or other cache.
// That might be done by implementing the 'Cache' interface
// and providing the one via 'Options'.
type inMemoryCache struct {
	mtx      sync.Mutex
	cache    map[string]*value
	capacity uint

	size  uint
	first *node
	last  *node
}

func (s *inMemoryCache) itemOnTop(v *value) {
	if v.item != s.last {
		if v.item == s.first {
			s.first = s.first.next
		} else {
			v.item.prev.next = v.item.next
			v.item.next.prev = v.item.prev
		}

		v.item.next = nil
		v.item.prev = s.last
		s.last.next = v.item
		s.last = v.item
	}
}

func (s *inMemoryCache) Put(key string, val []byte, _ time.Duration) error {
	if s.capacity < 1 {
		return errors.New("There is no possibility for an insertion. The capacity is 0.")
	}

	if len(key) == 0 {
		return errors.New("You can't use an empty string like a key.")
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	if v, ok := s.cache[key]; ok {
		v.val = val
		s.itemOnTop(v)
	} else {
		item := &node{
			key:  key,
			prev: s.last,
			next: nil,
		}

		v := &value{
			val:  val,
			item: item,
		}

		if s.first == nil {
			s.first = item
			s.last = item
		}

		if s.size < s.capacity {
			s.size++
		} else {
			delete(s.cache, s.first.key)

			s.first = s.first.next
			if s.first != nil {
				s.first.prev = nil
			} else {
				s.first = item
				s.last = item
			}

		}

		if s.last != item {
			s.last.next = item
			s.last = item
		}
		s.cache[key] = v
	}

	return nil
}

func (s *inMemoryCache) Get(key string) ([]byte, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if v, ok := s.cache[key]; ok {
		s.itemOnTop(v)

		return v.val, nil
	}

	return nil, nil
}

// Makes a built-in LRU in-memory cache implementation
// aimed for small projects, MVP, etc.
func NewInMemoryCache(capacity uint) Cache {
	return &inMemoryCache{
		capacity: capacity,
		cache:    make(map[string]*value, capacity),
		size:     0,
		first:    nil,
		last:     nil,
	}
}
