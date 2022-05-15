package limiter

import "sync"

type bucket struct {
	mu    *sync.Mutex
	limit int64
	hash  map[string]int64
}

func newBucket(limit int64) *bucket {
	return &bucket{
		mu:    &sync.Mutex{},
		limit: limit,
		hash:  make(map[string]int64),
	}
}

func (b *bucket) addKey(key string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.hash[key] >= b.limit {
		return false
	}

	b.hash[key]++

	return true
}

func (b *bucket) resetKey(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.hash[key] = 0
}

func (b *bucket) resetBucket() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.hash = make(map[string]int64)
}
