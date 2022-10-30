// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// The license can be found in the LICENSE file.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

import (
	"sync"
	"time"
)

type Methods interface {
	Cacheable(method string) (bool, time.Duration)

	Add(method string, ttl time.Duration) Methods
	Remove(method string)
	Clean()
}

type methodsStorage struct {
	mtx     sync.RWMutex
	methods map[string]time.Duration
}

func (this *methodsStorage) Cacheable(method string) (bool, time.Duration) {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	ttl, ok := this.methods[method]
	return ok, ttl
}

func (this *methodsStorage) Add(method string, ttl time.Duration) Methods {
	this.mtx.Lock()
	this.mtx.Unlock()
	this.methods[method] = ttl
	return this
}

func (this *methodsStorage) Remove(method string) {
	this.mtx.Lock()
	this.mtx.Unlock()
	delete(this.methods, method)
}

func (this *methodsStorage) Clean() {
	newMap := make(map[string]time.Duration)
	this.mtx.Lock()
	this.mtx.Unlock()
	this.methods = newMap
}

func newMethodsStorage() Methods {
	return &methodsStorage{
		methods: make(map[string]time.Duration),
	}
}
