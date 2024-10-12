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

	err := s.storage.AcceptReturn(ctx, req.GetClientId(), req.GetOrderId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &desc.AcceptReturnResponse{}, nil
}
