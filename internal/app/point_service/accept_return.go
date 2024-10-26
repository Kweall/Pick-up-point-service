package point_service

import (
	"context"
	desc "homework/pkg/point-service/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) AcceptReturn(ctx context.Context, req *desc.AcceptReturnRequest) (*desc.AcceptReturnResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	isReceived, isReturned, err := s.storage.CheckOrderStatus(ctx, req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check order status: %v", err)
	}

	if !isReceived {
		return nil, status.Errorf(codes.FailedPrecondition, "order %d has not been received yet", req.OrderId)
	}

	if isReturned {
		return nil, status.Errorf(codes.FailedPrecondition, "order %d has already been returned", req.OrderId)
	}

	err = s.storage.AcceptReturn(ctx, req.GetClientId(), req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	cacheKey := "all_returns"
	if err := s.cache.Invalidate(ctx, cacheKey); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete cache for returns")
	}
	return &desc.AcceptReturnResponse{}, nil
}
