// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

type zeroMetaFilter struct {
}

func (*zeroMetaFilter) Allowed(string, []string) bool {
	return true
}

// Makes a zero-filter implementation. It can be used
// if you don't have any rule for header filtering.
func NewZeroMetaFilter() MetaFilter {
	return &zeroMetaFilter{}
}
