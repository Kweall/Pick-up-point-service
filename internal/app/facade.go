package app

import (
	"context"
	"time"

	"homework/internal/storage/postgres"
)

//go:generate mkdir -p ./mocks
type Facade interface {
	AddOrder(ctx context.Context, orderID int64, clientID int64, createdAt time.Time, expiredAt time.Time, weight float64, price int64, packaging string, additional_film string) error
	DeleteOrder(ctx context.Context, orderID int64) error
	GetOrders(ctx context.Context, clientID int64) ([]*postgres.Order, error)
	GiveOrders(ctx context.Context, orderIDs []int64) error
	AcceptReturn(ctx context.Context, clientID, orderID int64) error
	CheckOrderStatus(ctx context.Context, orderID int64) (bool, bool, error)
	GetReturns(ctx context.Context) ([]*postgres.Order, error)
	GetOrdersByIDs(ctx context.Context, orderIDs []int64) ([]*postgres.Order, error)
}

type storageFacade struct {
	txManager    postgres.TransactionManager
	pgRepository *postgres.PgRepository
}

func NewStorageFacade(
	txManager postgres.TransactionManager,
	pgRepository *postgres.PgRepository,
) *storageFacade {
	return &storageFacade{
		txManager:    txManager,
		pgRepository: pgRepository,
	}
}

func (s *storageFacade) AddOrder(ctx context.Context, orderID int64, clientID int64, createdAt time.Time, expiredAt time.Time, weight float64, price int64, packaging string, additional_film string) error {
	return s.txManager.RunSerializable(ctx, func(ctxTx context.Context) error {
		err := s.pgRepository.AddOrderHistory(ctx, orderID)
		if err != nil {
			return err
		}

		err = s.pgRepository.AddOrder(ctx, orderID, clientID, createdAt, expiredAt, weight, price, packaging, additional_film)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *storageFacade) DeleteOrder(ctx context.Context, orderID int64) error {
	return s.pgRepository.DeleteOrder(ctx, orderID)
}

func (s *storageFacade) GetOrders(ctx context.Context, userID int64) ([]*postgres.Order, error) {
	var contacts []*postgres.Order
	err := s.txManager.RunReadUncommitted(ctx, func(ctxTx context.Context) error {
		c, err := s.pgRepository.GetOrders(ctx, userID)
		if err != nil {
			return err
		}

		contacts = c
		return nil
	})

	return contacts, err
}

func (s *storageFacade) GiveOrders(ctx context.Context, orderIDs []int64) error {
	return s.pgRepository.GiveOrders(ctx, orderIDs)
}

func (s *storageFacade) AcceptReturn(ctx context.Context, clientID, orderID int64) error {
	return s.pgRepository.AcceptReturn(ctx, clientID, orderID)
}

func (s *storageFacade) CheckOrderStatus(ctx context.Context, orderID int64) (bool, bool, error) {
	return s.pgRepository.CheckOrderStatus(ctx, orderID)
}

func (s *storageFacade) GetReturns(ctx context.Context) ([]*postgres.Order, error) {
	return s.pgRepository.GetReturns(ctx)
}

func (s *storageFacade) GetOrdersByIDs(ctx context.Context, orderIDs []int64) ([]*postgres.Order, error) {
	return s.pgRepository.GetOrdersByIDs(ctx, orderIDs)
}
