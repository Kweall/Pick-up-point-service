package point_service

import (
	"context"
	desc "homework/pkg/point-service/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const returnsPerPage = 5

func (s *Implementation) GetReturns(ctx context.Context, req *desc.GetReturnsRequest) (*desc.GetReturnsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	orders, err := s.storage.GetReturns(ctx)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if len(orders) == 0 {
		return &desc.GetReturnsResponse{Orders: nil}, nil
	}

	start := (int(req.Page) - 1) * returnsPerPage
	if start >= len(orders) {
		return nil, status.Error(codes.InvalidArgument, "page number exceeds the available range")
	}

	end := start + returnsPerPage
	if end > len(orders) {
		end = len(orders)
	}

	res := &desc.GetReturnsResponse{
		Orders: make([]*desc.Order, 0, len(orders)),
	}
	for _, order := range orders[start:end] {
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
