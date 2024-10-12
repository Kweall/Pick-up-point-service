package point_service

import (
	"context"
	desc "homework/pkg/point-service/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) DeleteOrder(ctx context.Context, req *desc.DeleteOrderRequest) (*desc.DeleteOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	err := s.storage.DeleteOrder(ctx, req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &desc.DeleteOrderResponse{}, nil
}
