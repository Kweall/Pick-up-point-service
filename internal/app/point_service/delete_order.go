package point_service

import (
	"context"
	"fmt"
	desc "homework/pkg/point-service/v1"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) DeleteOrder(ctx context.Context, req *desc.DeleteOrderRequest) (*desc.DeleteOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	clientID, err := s.storage.DeleteOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	cacheKey := fmt.Sprintf("orders:client:%d", clientID)
	if err := s.cache.Invalidate(ctx, cacheKey); err != nil {
		log.Printf("failed to invalidate cache for clientId %d: %v", clientID, err)
	}
	return &desc.DeleteOrderResponse{}, nil
}
