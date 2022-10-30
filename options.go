// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// The license can be found in the LICENSE file.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

type Options struct {
	Switch     Switch
	Methods    Methods
	MetaFilter MetaFilter
	Hash       Hash
	Cache      Cache
}

func NewOptions() *Options {
	opts := Options{
		Switch:     newSwitch(),
		Methods:    newMethodsStorage(),
		MetaFilter: NewZeroMetaFilter(),
		Hash:       newDefaultHash(),
		Cache:      NewInMemoryCache(1024),
	}

	return &opts
}
