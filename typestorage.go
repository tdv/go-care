// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

import (
	"errors"
	"fmt"
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

func (this *typeStorageImpl) Put(key string, val interface{}) (err error) {
	if len(key) == 0 {
		err = errors.New("Key must not be empty.")
		return
	}

	typ := reflect.TypeOf(val)
	if typ == nil {
		err = errors.New("Failed to get type of the value.")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint("The value is not an interface. ", r))
		}
	}()

	elem := typ.Elem()

	this.mtx.Lock()
	defer this.mtx.Unlock()
	this.types[key] = elem

	return
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
