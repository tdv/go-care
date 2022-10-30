// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

import "time"

// A 'Cache' represents a common interface for responses caching.
// It can be implemented for many different caches like
// Redis, Memcached, etc,
type Cache interface {
	Put(key string, val []byte, ttl time.Duration) error
	Get(key string) ([]byte, error)
}
