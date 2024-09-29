package postgres

import (
	"context"
	"homework/internal/storage/postgres"
	"testing"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestE2E_GiveOrders(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	require.NoError(t, err)
	defer pool.Close()

	txManager := postgres.NewTxManager(pool)
	repo := postgres.NewPgRepository(txManager)

	err = repo.AddOrder(ctx, 10, 100, time.Now(), time.Now().Add(24*time.Hour), 1.5, 500, "Box", "None")
	require.NoError(t, err)

	err = repo.GiveOrders(ctx, []int64{10})
	require.NoError(t, err)

	var receivedAt time.Time
	err = pool.QueryRow(ctx, "SELECT received_at FROM orders WHERE order_id = $1", 10).Scan(&receivedAt)
	require.NoError(t, err)
	require.False(t, receivedAt.IsZero())
}
