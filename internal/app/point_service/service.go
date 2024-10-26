package point_service

import (
	"homework/internal/app"
	"homework/internal/cache"
	desc "homework/pkg/point-service/v1"
)

type Implementation struct {
	storage app.Facade
	cache   *cache.RedisCache

	desc.UnimplementedPointServiceServer
}

func NewImplementation(storage app.Facade, cache *cache.RedisCache) *Implementation {
	return &Implementation{
		storage: storage,
		cache:   cache}
}
