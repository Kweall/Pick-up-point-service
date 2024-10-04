package app

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

type AcceptReturnRequest struct {
	OrderID    int64      `json:"OrderID"`
	ClientID   int64      `json:"ClientID"`
	ReceivedAt *time.Time `json:"ReceivedAt"`
}

type AcceptReturnResponse struct {
}

const countOfArgumentstoReturn = 2

func (s *Service) AcceptReturn(ctx context.Context, req *AcceptReturnRequest, parts []string) (*AcceptReturnResponse, error) {
	var err error
	if len(parts)-1 != countOfArgumentstoReturn {
		return nil, fmt.Errorf("should be 2 arguments: clientID (int), orderID (int)")
	}

	req.ClientID, err = strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("clientID is incorrect")
	}

	req.OrderID, err = strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("orderID is incorrect")
	}

	isReceived, isReturned, err := s.storage.CheckOrderStatus(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check order status: %v", err)
	}

	if !isReceived {
		return nil, fmt.Errorf("order %d has not been received yet", req.OrderID)
	}

	if isReturned {
		return nil, fmt.Errorf("order %d has already been returned", req.OrderID)
	}
	err = s.storage.AcceptReturn(ctx, req.ClientID, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	fmt.Printf("Return accepted for Order ID: %d, ClientID: %d\n", req.OrderID, req.ClientID)

	return &AcceptReturnResponse{}, nil
}
