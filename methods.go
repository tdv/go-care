package care

import (
	"sync"
)

type Methods interface {
	Exists(method string) bool

	Add(method string) Methods
	Set(methods ...string)
	Remove(method string)
	Clean()
}

type methodsStorage struct {
	mtx     sync.RWMutex
	methods map[string]struct{}
}

func (this *methodsStorage) Exists(method string) bool {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	_, ok := this.methods[method]
	return ok
}

func (this *methodsStorage) Add(method string) Methods {
	this.mtx.Lock()
	this.mtx.Unlock()
	this.methods[method] = struct{}{}
	return this
}

func (this *methodsStorage) Set(methods ...string) {
	newMap := make(map[string]struct{})
	for _, method := range methods {
		newMap[method] = struct{}{}
	}
	this.mtx.Lock()
	this.mtx.Unlock()
	this.methods = newMap
}

func (this *methodsStorage) Remove(method string) {
	this.mtx.Lock()
	this.mtx.Unlock()
	delete(this.methods, method)
}

func (this *methodsStorage) Clean() {
	newMap := make(map[string]struct{})
	this.mtx.Lock()
	this.mtx.Unlock()
	this.methods = newMap
}

func newMethodsStorage() Methods {
	return &methodsStorage{
		methods: make(map[string]struct{}),
	}
}
