package app

import (
	"context"
	"fmt"
	"homework/internal/storage/postgres"
	"sort"
	"strconv"
)

type GetOrdersRequest struct {
	ClientID int64 `json:"ClientID"`
	Limit    int   `json:"Limit,omitempty"`
}

type GetOrdersResponse struct {
	Orders []*postgres.Order `json:"orders"`
}

const countOfArgumentsToGetOrders = 2

func (s *Service) GetOrders(ctx context.Context, req *GetOrdersRequest, parts []string) (*GetOrdersResponse, error) {
	if len(parts)-1 < countOfArgumentsToGetOrders-1 {
		return nil, fmt.Errorf("should be at least 1 argument: clientID (int)")
	} else if len(parts)-1 > countOfArgumentsToGetOrders {
		return nil, fmt.Errorf("should be maximum 2 arguments: clientID (int) and count of orders you want to get")
	}

	var err error
	req.ClientID, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("clientID is incorrect")
	}

	req.Limit = -1 // -1 будет означать, что ограничение не задано
	if len(parts)-1 == countOfArgumentsToGetOrders {
		req.Limit, err = strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("limit is incorrect")
		}
	}

	orders, err := s.storage.GetOrders(ctx, req.ClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %v", err)
	}

	if len(orders) == 0 {
		return nil, fmt.Errorf("no orders found for clientID: %v", req.ClientID)
	}

	sort.Slice(orders, func(i, j int) bool {
		if orders[i].CreatedAt == nil && orders[j].CreatedAt == nil {
			return false
		}
		if orders[i].CreatedAt == nil {
			return false
		}
		if orders[j].CreatedAt == nil {
			return true
		}
		return orders[i].CreatedAt.After(*orders[j].CreatedAt)
	})

	if req.Limit > 0 && req.Limit < len(orders) {
		orders = orders[:req.Limit]
	}

	resp := &GetOrdersResponse{
		Orders: make([]*postgres.Order, 0, len(orders)),
	}

	for _, order := range orders {
		resp.Orders = append(resp.Orders, &postgres.Order{
			OrderID:  order.OrderID,
			ClientID: order.ClientID,
		})
	}

	for _, order := range resp.Orders {
		fmt.Printf("Order ID: %d, Client ID: %d\n", order.OrderID, order.ClientID)
	}

	return resp, nil
}
