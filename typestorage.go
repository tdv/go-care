package care

import (
	"errors"
	"reflect"
	"sync"
)

type typeStorage interface {
	Put(key string, val interface{}) error
	Get(key string) (reflect.Type, bool)
}

type typeStorageImpl struct {
	mtx   sync.RWMutex
	types map[string]reflect.Type
}

func (this *typeStorageImpl) Put(key string, val interface{}) error {
	if len(key) == 0 {
		return errors.New("Key must not be empty.")
	}

	typ := reflect.TypeOf(val)
	if typ == nil {
		return errors.New("Failed to get type of the value.")
	}

	var err error = nil

	defer func() {
		if recover() != nil {
			err = errors.New("The value is not an interface.")
		}
	}()

	elem := typ.Elem()

	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.types[key] = elem

	return err
}

func (this *typeStorageImpl) Get(key string) (reflect.Type, bool) {
	this.mtx.RLock()
	defer this.mtx.RUnlock()

	typ, ok := this.types[key]

	return typ, ok
}

func newTypeStorage() typeStorage {
	return &typeStorageImpl{
		types: make(map[string]reflect.Type),
	}
}
