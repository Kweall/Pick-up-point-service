package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/prometheus/client_golang/prometheus"
)

type RedisCache struct {
	client *redis.Client
	ttl    time.Duration
}

var (
	RedisCacheHits   = prometheus.NewCounter(prometheus.CounterOpts{Name: "redis_cache_hits_total", Help: "Total number of Redis cache hits"})
	RedisCacheMisses = prometheus.NewCounter(prometheus.CounterOpts{Name: "redis_cache_misses_total", Help: "Total number of Redis cache misses"})
)

func init() {
	prometheus.MustRegister(RedisCacheHits)
	prometheus.MustRegister(RedisCacheMisses)
}

func NewRedisCache(addr, password string, db int, ttl time.Duration) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisCache{client: client, ttl: ttl}
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		RedisCacheMisses.Inc()
		return "", nil
	} else if err != nil {
		return "", err
	}
	RedisCacheHits.Inc()
	return value, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value string) error {
	return c.client.Set(ctx, key, value, c.ttl).Err()
}

func (c *RedisCache) Invalidate(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
