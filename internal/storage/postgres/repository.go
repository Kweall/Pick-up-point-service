package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx"
)

var (
	ErrorNotFoundOrder = errors.New("order not found")
)

type PgRepository struct {
	txManager TransactionManager
}

func NewPgRepository(txManager TransactionManager) *PgRepository {
	return &PgRepository{
		txManager: txManager,
	}
}

func (r *PgRepository) AddOrder(ctx context.Context, req *Order) error {
	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `
        INSERT INTO orders (order_id, client_id, created_at, expired_at, weight, price, packaging, additional_film)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `, req.OrderID, req.ClientID, req.CreatedAt, req.ExpiredAt, req.Weight, req.Price, req.Packaging, req.AdditionalFilm)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgRepository) AddOrderHistory(ctx context.Context, orderID int64) error {
	tx := r.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `
		insert into orders_history(order_id) values ($1)
	`, orderID)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgRepository) DeleteOrder(ctx context.Context, orderID int64) (clientID int64, err error) {
	tx := r.txManager.GetQueryEngine(ctx)

	err = tx.QueryRow(ctx, `
        select client_id from orders where order_id = $1
    `, orderID).Scan(&clientID)
	if err != nil {
		return 0, fmt.Errorf("failed to query client ID for order ID %d: %w", orderID, err)
	}

	result, err := tx.Exec(ctx, `
		delete from orders where order_id = $1
	`, orderID)
	if err != nil {
		return 0, err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return 0, ErrorNotFoundOrder
	}
	fmt.Printf("Order %v has been deleted from database\n", orderID)
	return clientID, nil
}

func (r *PgRepository) GetOrders(ctx context.Context, clientID int64) ([]*Order, error) {
	var orders []*Order
	tx := r.txManager.GetQueryEngine(ctx)
	err := pgxscan.Select(ctx, tx, &orders, `
		select client_id, order_id from orders where client_id = $1
	`, clientID)

	return orders, err
}

func (r *PgRepository) GiveOrders(ctx context.Context, orderIDs []int64) error {
	tx := r.txManager.GetQueryEngine(ctx)
	_, err := tx.Exec(ctx, `
        UPDATE orders SET received_at = $1 WHERE order_id = ANY($2)
    `, time.Now(), orderIDs)
	return err
}

func (r *PgRepository) UpdateReceivedAt(ctx context.Context, orderIDs []int64, receivedAt time.Time) error {
	tx := r.txManager.GetQueryEngine(ctx)

	orderIDParams := make([]string, len(orderIDs))
	for i := range orderIDs {
		orderIDParams[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
		UPDATE orders SET received_at = $%d WHERE order_id IN (%s)
	`, len(orderIDs)+1, strings.Join(orderIDParams, ", "))

	args := make([]interface{}, len(orderIDs)+1)
	for i, id := range orderIDs {
		args[i] = id
	}
	args[len(orderIDs)] = receivedAt

	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *PgRepository) GetOrdersByIDs(ctx context.Context, orderIDs []int64) ([]*Order, error) {
	var orders []*Order

	orderIDParams := make([]string, len(orderIDs))
	for i := range orderIDs {
		orderIDParams[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf(`
        SELECT order_id, client_id, received_at FROM orders WHERE order_id IN (%s)
    `, strings.Join(orderIDParams, ", "))

	args := make([]interface{}, len(orderIDs))
	for i, id := range orderIDs {
		args[i] = id
	}

	err := pgxscan.Select(ctx, r.txManager.GetQueryEngine(ctx), &orders, query, args...)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PgRepository) AcceptReturn(ctx context.Context, clientID, orderID int64) error {
	tx := r.txManager.GetQueryEngine(ctx)

	_, err := tx.Exec(ctx, `
        UPDATE orders SET returned_at = $1 WHERE client_id = $2 AND order_id = $3
    `, time.Now(), clientID, orderID)

	if err != nil {
		return fmt.Errorf("failed to accept return: %w", err)
	}

	return nil
}

func (r *PgRepository) CheckOrderStatus(ctx context.Context, orderID int64) (bool, bool, error) {
	var receivedAt *time.Time
	var returnedAt *time.Time

	err := r.txManager.GetQueryEngine(ctx).QueryRow(ctx, `
		SELECT received_at, returned_at
		FROM orders
		WHERE order_id = $1
	`, orderID).Scan(&receivedAt, &returnedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, false, fmt.Errorf("order not found: %d", orderID)
		}
		return false, false, fmt.Errorf("failed to check order status: %w", err)
	}
	// Определяем, получен ли заказ и был ли он возвращен
	isReceived := receivedAt != nil
	isReturned := returnedAt != nil

	return isReceived, isReturned, nil
}

func (r *PgRepository) GetReturns(ctx context.Context) ([]*Order, error) {
	var orders []*Order

	err := pgxscan.Select(ctx, r.txManager.GetQueryEngine(ctx), &orders, `
		SELECT order_id, client_id, returned_at
		FROM orders
		WHERE returned_at IS NOT NULL
	`)
	if err != nil {
		return nil, err
	}

	return orders, nil
}
