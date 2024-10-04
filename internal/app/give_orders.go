package app

import (
	"context"
	"fmt"
	"strconv"
)

type GiveOrderRequest struct {
	OrderIDs []int64 `json:"OrderIDs"`
}

type GiveOrderResponse struct {
	UpdatedCount int64 `json:"updated_count"`
}

func (s *Service) GiveOrders(ctx context.Context, req *GiveOrderRequest, parts []string) (*GiveOrderResponse, error) {
	var err error
	if len(parts)-1 < 1 {
		return nil, fmt.Errorf("should be at least 1 argument: list of orderID's (int) separated by space")
	}

	for _, part := range parts[1:] {
		orderID, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("orderID is incorrect")
		}
		req.OrderIDs = append(req.OrderIDs, orderID)
	}

	orders, err := s.storage.GetOrdersByIDs(ctx, req.OrderIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("no orders found")
	}

	firstClientID := orders[0].ClientID
	for _, order := range orders {
		if order.ClientID != firstClientID {
			return nil, fmt.Errorf("all orders must belong to the same client")
		}
		if order.ReceivedAt != nil {
			return nil, fmt.Errorf("order %d has already been received", order.OrderID)
		}
	}

	err = s.storage.GiveOrders(ctx, req.OrderIDs)
	if err != nil {
		return nil, err
	}
	if len(req.OrderIDs) == 1 {
		fmt.Printf("The order successfully given to client\n")
	} else {
		fmt.Printf("These orders successfully given to client\n")
	}
	return &GiveOrderResponse{UpdatedCount: int64(len(req.OrderIDs))}, nil
}
