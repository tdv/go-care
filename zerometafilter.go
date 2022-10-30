// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// The license can be found in the LICENSE file.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

type zeroMetaFilter struct {
}

func (*zeroMetaFilter) Allowed(string, []string) bool {
	return true
}

func NewZeroMetaFilter() MetaFilter {
	return &zeroMetaFilter{}
}
