package point_service

import (
	"context"
	desc "homework/pkg/point-service/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GetOrders(ctx context.Context, req *desc.GetOrdersRequest) (*desc.GetOrdersResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orders, err := s.storage.GetOrders(ctx, req.GetClientId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	res := &desc.GetOrdersResponse{
		Orders: make([]*desc.Order, 0, len(orders)),
	}
	if len(orders) == 0 {
		return nil, status.Error(codes.NotFound, "client not found")
	}
	for _, order := range orders {
		res.Orders = append(res.Orders, &desc.Order{
			OrderId:  order.OrderID,
			ClientId: order.ClientID,
			// CreatedAt:      order.CreatedAt,
			// ExpiredAt:      order.ExpiredAt,
			// Weight:         float32(order.Weight),
			// Price:          order.Price,
			// Packaging:      order.Packaging,
			// AdditionalFilm: order.AdditionalFilm,
		})
	}
	return res, nil
}
