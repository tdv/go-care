// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

// Headers is a pool for filtering.
type Headers struct {
	// Headers for the key computation.
	Allowed []string
	// Omitted headers from the key computation.
	Disallowed []string
}

// MetaFilter represents an interface for filtering the headers
// before including the once in the key computation.
//
// It can be useful if you need to pick up only a few header
// which have to be included in to the key, making the one
// more unique and filter the noise. For instance, request-id,
// trace-id, etc are the noise, meanwhile jwt-token (and others
// according your app logic) is an important header.
//
// Having implemented your own version, you can control the headers
// which will be involved in the key computation process.
type MetaFilter interface {
	// Allowed - returns true allowing to include the header in
	// the key computation, otherwise returns false.
	Allowed(key string, val []string) bool
}

type metaFilter struct {
	allowedHeaders    map[string]struct{}
	disallowedHeaders map[string]struct{}
}

func (s *metaFilter) Allowed(key string, _ []string) bool {
	if _, ok := s.disallowedHeaders[key]; ok {
		return false
	}

	_, ok := s.allowedHeaders[key]
	return ok
}

//goland:noinspection ALL
func NewMetaFilter(headers Headers) MetaFilter {
	inst := metaFilter{
		allowedHeaders:    make(map[string]struct{}),
		disallowedHeaders: make(map[string]struct{}),
	}

	if len(headers.Allowed) > 0 {
		for _, h := range headers.Allowed {
			inst.allowedHeaders[h] = struct{}{}
		}
	}

	if len(headers.Disallowed) > 0 {
		for _, h := range headers.Disallowed {
			inst.disallowedHeaders[h] = struct{}{}
		}
	}

	return &inst
}
