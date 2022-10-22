package care

type Cache interface {
	Put(key string, val []byte) error
	Get(key string) ([]byte, error)
}
