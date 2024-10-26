package app

import (
	"log"
	"net/http"
	"sync"

	"homework/internal/cache"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ordersGiven = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "orders_given_total",
		Help: "Total number of orders given out",
	})

	once sync.Once
)

func InitMetrics() {
	once.Do(func() {
		prometheus.MustRegister(
			cache.CacheHits,
			cache.CacheMisses,
			cache.RedisCacheHits,
			cache.RedisCacheMisses,
			ordersGiven,
		)
	})
}

func IncrementOrdersGiven() {
	ordersGiven.Inc()
}

func StartMetricsEndpoint() {
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":7002", nil); err != nil {
		log.Fatalf("Failed to start metrics endpoint: %v", err)
	}
}
