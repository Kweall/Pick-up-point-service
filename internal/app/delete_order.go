package app

import (
	"context"
	"fmt"
	"strconv"
)

type DeleteOrderRequest struct {
	OrderID int64 `json:"OrderID"`
}

type DeleteOrderResponse struct {
}

const countOfArgumentsToDelete = 1

func (s *Service) DeleteOrder(ctx context.Context, req *DeleteOrderRequest, parts []string) (*DeleteOrderResponse, error) {
	var err error
	if len(parts)-1 != countOfArgumentsToDelete {
		return nil, fmt.Errorf("should be 1 argument: orderID (int)")
	}
	req.OrderID, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("orderID is incorrect")
	}
	err = s.storage.DeleteOrder(ctx, req.OrderID)
	return nil, err
}
