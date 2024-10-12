package point_service

import (
	"homework/internal/app"
	desc "homework/pkg/point-service/v1"
)

type Implementation struct {
	storage app.Facade

	desc.UnimplementedPointServiceServer
}

func NewImplementation(storage app.Facade) *Implementation {
	return &Implementation{storage: storage}
}
