package point_service

import (
	"context"
	desc "homework/pkg/point-service/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GiveOrder(ctx context.Context, req *desc.GiveOrderRequest) (*desc.GiveOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.storage.GiveOrders(ctx, req.GetOrderIds())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &desc.GiveOrderResponse{}, nil
}
