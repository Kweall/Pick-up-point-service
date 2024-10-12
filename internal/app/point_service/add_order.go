package point_service

import (
	"context"
	"homework/internal/storage/postgres"
	desc "homework/pkg/point-service/v1"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Вспомогательная функция для преобразования time.Time в *time.Time
func TimeToPtr(t time.Time) *time.Time {
	return &t
}

func (s *Implementation) AddOrder(ctx context.Context, req *desc.AddOrderRequest) (*desc.AddOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.storage.AddOrder(ctx, &postgres.Order{
		OrderID:        req.OrderId,
		ClientID:       req.ClientId,
		CreatedAt:      TimeToPtr(req.CreatedAt.AsTime()), // Преобразование в *time.Time
		ExpiredAt:      TimeToPtr(req.ExpiredAt.AsTime()), // Преобразование в *time.Time
		Weight:         float64(req.Weight),
		Price:          req.Price,
		Packaging:      req.Packaging,
		AdditionalFilm: *req.AdditionalFilm,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &desc.AddOrderResponse{}, nil
}
