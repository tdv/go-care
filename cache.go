package care

import "time"

type Cache interface {
	Put(key string, val []byte, ttl time.Duration) error
	Get(key string) ([]byte, error)
}
