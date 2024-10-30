package point_service

import (
	"context"
	"fmt"
	"homework/internal/app"
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

const timeLayout = "02.01.2006"

func (s *Implementation) AddOrder(ctx context.Context, req *desc.AddOrderRequest) (*desc.AddOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	createdAt := time.Now().Truncate(time.Minute)
	expiredAt, err := time.Parse(timeLayout, req.ExpiredAt.AsTime().Format(timeLayout))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid expired_at format, should be YYYY-MM-DDT00:00:00Z")
	}

	if expiredAt.Before(createdAt) {
		return nil, status.Error(codes.InvalidArgument, "expiration date must be today or later")
	}

	packaging, err := app.GetPackaging(req.Packaging)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to get packaging")
	}
	if !packaging.CheckWeight(float64(req.Weight)) {
		return nil, status.Error(codes.FailedPrecondition, "weight is too big for the selected packaging")
	}
	req.Price += packaging.GetPrice()
	if *req.AdditionalFilm {
		if req.Packaging != "film" {
			req.Price += 1
		} else {
			return nil, status.Error(codes.FailedPrecondition, "can't add additional film to film")
		}
	}

	err = s.storage.AddOrder(ctx, &postgres.Order{
		OrderID:        req.OrderId,
		ClientID:       req.ClientId,
		CreatedAt:      &createdAt,
		ExpiredAt:      TimeToPtr(req.ExpiredAt.AsTime()),
		Weight:         float64(req.Weight),
		Price:          req.Price,
		Packaging:      req.Packaging,
		AdditionalFilm: *req.AdditionalFilm,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	cacheKey := fmt.Sprintf("orders:client:%d", req.ClientId)
	if err := s.cache.Invalidate(ctx, cacheKey); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to delete cache for client %d: %v", req.ClientId, err))
	}
	return &desc.AddOrderResponse{}, nil
}
