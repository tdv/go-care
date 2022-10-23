package care

import (
	"errors"
	"sync"
)

type node struct {
	key  string
	prev *node
	next *node
}

type value struct {
	item *node
	val  []byte
}

type inMemoryCache struct {
	mtx      sync.Mutex
	cache    map[string]*value
	capacity uint

	size  uint
	first *node
	last  *node
}

func (this *inMemoryCache) itemOnTop(v *value) {
	if v.item != this.last {
		if v.item == this.first {
			this.first = this.first.next
		} else {
			v.item.prev.next = v.item.next
			v.item.next.prev = v.item.prev
		}

		v.item.next = nil
		v.item.prev = this.last
		this.last.next = v.item
		this.last = v.item
	}
}

func (this *inMemoryCache) Put(key string, val []byte) error {
	if this.capacity < 1 {
		return errors.New("There is no possibility for an insertion. The capacity is 0.")
	}

	if len(key) == 0 {
		return errors.New("You can't use an empty string like a key.")
	}

	this.mtx.Lock()
	defer this.mtx.Unlock()

	if v, ok := this.cache[key]; ok {
		v.val = val
		this.itemOnTop(v)
	} else {
		item := &node{
			key:  key,
			prev: this.last,
			next: nil,
		}

		v := &value{
			val:  val,
			item: item,
		}

		if this.first == nil {
			this.first = item
			this.last = item
		}

		if this.size < this.capacity {
			this.size++
		} else {
			delete(this.cache, this.first.key)

			this.first = this.first.next
			if this.first != nil {
				this.first.prev = nil
			} else {
				this.first = item
				this.last = item
			}

		}

		if this.last != item {
			this.last.next = item
			this.last = item
		}
		this.cache[key] = v
	}

	return nil
}

func (this *inMemoryCache) Get(key string) ([]byte, error) {
	this.mtx.Lock()
	defer this.mtx.Unlock()

	if v, ok := this.cache[key]; ok {
		this.itemOnTop(v)

		return v.val, nil
	}

	return nil, nil
}

func NewInMemoryCache(capacity uint) Cache {
	return &inMemoryCache{
		capacity: capacity,
		cache:    make(map[string]*value, capacity),
		size:     0,
		first:    nil,
		last:     nil,
	}
}
