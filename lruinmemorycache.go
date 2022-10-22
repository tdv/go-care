package care

type lruInMemoryCache struct {
	size uint
}

func (this *lruInMemoryCache) Put(key string, val []byte) error {
	return nil
}

func (this *lruInMemoryCache) Get(key string) ([]byte, error) {
	return nil, nil
}

func NewInMemoryLRUCache(size uint) Cache {
	return &lruInMemoryCache{
		size: size,
	}
}
