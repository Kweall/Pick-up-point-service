package point_service

import (
	"context"
	"encoding/json"
	"fmt"
	desc "homework/pkg/point-service/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Implementation) GetOrders(ctx context.Context, req *desc.GetOrdersRequest) (*desc.GetOrdersResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	cacheKey := fmt.Sprintf("orders:client:%d", req.GetClientId())

	cachedData, err := s.cache.Get(ctx, cacheKey)
	if err == nil && cachedData != "" {
		var cachedResponse desc.GetOrdersResponse
		if err := json.Unmarshal([]byte(cachedData), &cachedResponse); err == nil {
			return &cachedResponse, nil
		}
	}

	orders, err := s.storage.GetOrders(ctx, req.GetClientId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	res := &desc.GetOrdersResponse{
		Orders: make([]*desc.Order, 0, len(orders)),
	}
	for _, order := range orders {
		res.Orders = append(res.Orders, &desc.Order{
			OrderId:  order.OrderID,
			ClientId: order.ClientID,
		})
	}
	dataToCache, err := json.Marshal(res)
	if err == nil {
		if err := s.cache.Set(ctx, cacheKey, string(dataToCache)); err != nil {
			return nil, status.Error(codes.Aborted, err.Error())
		}
	}
	return res, nil
}
