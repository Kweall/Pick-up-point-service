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

	orders, err := s.storage.GetOrdersByIDs(ctx, req.GetOrderIds())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get orders: %v", err)
	}

	if len(orders) != len(req.GetOrderIds()) {
		return nil, status.Error(codes.NotFound, "one or more orders not found")
	}

	firstClientID := orders[0].ClientID
	for _, order := range orders {
		if order.ClientID != firstClientID {
			return nil, status.Error(codes.InvalidArgument, "all orders must belong to the same client")
		}
		if order.ReceivedAt != nil {
			return nil, status.Errorf(codes.FailedPrecondition, "order %d has already been received", order.OrderID)
		}
	}

	err = s.storage.GiveOrders(ctx, req.GetOrderIds())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &desc.GiveOrderResponse{}, nil
}
