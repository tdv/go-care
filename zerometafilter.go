// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

type zeroMetaFilter struct {
	allowed bool
}

func (s *zeroMetaFilter) Allowed(string, []string) bool {
	return s.allowed
}

// Makes a zero-filter implementation. It can be used
// if you don't have any rule for header filtering.
//
// allowed - defines common behaviour for any header
//
//	true - allows all header to be included in the key computation.
//	false - disallows all header to be included in the key computation.
func NewZeroMetaFilter(allowed bool) MetaFilter {
	return &zeroMetaFilter{
		allowed: allowed,
	}
}
