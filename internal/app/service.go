package app

import (
	"context"

	"homework/internal/storage/postgres"
)

//go:generate mkdir -p ./mocks
type Facade interface {
	AddOrder(ctx context.Context, req *postgres.Order) error
	DeleteOrder(ctx context.Context, orderID int64) (int64, error)
	GetOrders(ctx context.Context, clientID int64) ([]*postgres.Order, error)
	GiveOrders(ctx context.Context, orderIDs []int64) error
	AcceptReturn(ctx context.Context, clientID, orderID int64) error
	CheckOrderStatus(ctx context.Context, orderID int64) (bool, bool, error)
	GetReturns(ctx context.Context) ([]*postgres.Order, error)
	GetOrdersByIDs(ctx context.Context, orderIDs []int64) ([]*postgres.Order, error)
}

type Service struct {
	storage Facade
}

func NewService(storage Facade) *Service {
	return &Service{storage: storage}
}
