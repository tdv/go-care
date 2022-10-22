// Copyright (c) 2022 Dmitry Tkachenko (tkachenkodmitryv@gmail.com)
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package care

import (
	"context"
	"time"
)

// Cache represents a common interface for responses caching.
// It can be implemented for many caches like Redis, Memcached, etc,
type Cache interface {
	// Put data into the cache by key.
	Put(ctx context.Context, key string, val []byte, ttl time.Duration) error
	// Get data from the cache by key.
	Get(ctx context.Context, key string) ([]byte, error)
}
