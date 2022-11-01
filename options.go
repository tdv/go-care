// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

// Memoization options. There are some options you can redefine
// by your own implementation in order to have more flexibility.
// 'Options' gives enough room space to do that.
type Options struct {
	// The memoization feature is turned on/off.
	Switch Switch
	// A pool of methods for memorizing responses.
	Methods Methods
	// Header filter for including the ones in the key computation.
	MetaFilter MetaFilter
	// Used for the key computation.
	Hash Hash
	// A cache
	Cache Cache
}

// Makes a default options set, having filled
// all items by thread-safe implementations.
func NewOptions() *Options {
	opts := Options{
		Switch:     newSwitch(),
		Methods:    newMethodsStorage(),
		MetaFilter: NewZeroMetaFilter(true),
		Hash:       newDefaultHash(),
		Cache:      NewInMemoryCache(1024),
	}

	return &opts
}
