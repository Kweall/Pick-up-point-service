package cache

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type MemoryCache[K comparable, V any] struct {
	mu       sync.Mutex
	data     map[K]cacheItem[V]
	capacity int
	ttl      time.Duration
}

type cacheItem[V any] struct {
	value     V
	timestamp time.Time
}

var (
	CacheHits   = prometheus.NewCounter(prometheus.CounterOpts{Name: "memory_cache_hits_total", Help: "Total number of memory cache hits"})
	CacheMisses = prometheus.NewCounter(prometheus.CounterOpts{Name: "memory_cache_misses_total", Help: "Total number of memory cache misses"})
)

func init() {
	prometheus.MustRegister(CacheHits)
	prometheus.MustRegister(CacheMisses)
}

func NewMemoryCache[K comparable, V any](capacity int, ttl time.Duration) *MemoryCache[K, V] {
	return &MemoryCache[K, V]{
		data:     make(map[K]cacheItem[V]),
		capacity: capacity,
		ttl:      ttl,
	}
}

func (c *MemoryCache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, found := c.data[key]
	if !found || time.Since(item.timestamp) > c.ttl {
		if found {
			delete(c.data, key)
		}
		CacheMisses.Inc()
		var zero V
		return zero, false
	}

	CacheHits.Inc()
	return item.value, true
}

func (c *MemoryCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.data) >= c.capacity {
		c.evict()
	}

	c.data[key] = cacheItem[V]{value: value, timestamp: time.Now()}
}

func (c *MemoryCache[K, V]) evict() {
	var oldestKey K
	oldestTime := time.Now()

	for k, item := range c.data {
		if item.timestamp.Before(oldestTime) {
			oldestTime = item.timestamp
			oldestKey = k
		}
	}

	delete(c.data, oldestKey)
}
