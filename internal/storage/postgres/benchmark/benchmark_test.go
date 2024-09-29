package postgres

import (
	"context"
	"testing"
	"time"

	"homework/internal/storage/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func setupPostgresStorage(b *testing.B) (*postgres.PgRepository, *pgxpool.Pool, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	// Устанавливаем соединение с PostgreSQL
	pool, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	require.NoError(b, err)

	// Создаем менеджер транзакций и репозиторий
	txManager := postgres.NewTxManager(pool)
	repo := postgres.NewPgRepository(txManager)

	// Создаем таблицы перед тестом
	_, err = pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			order_id BIGINT PRIMARY KEY,
			client_id BIGINT,
			created_at TIMESTAMPTZ,
			expired_at TIMESTAMPTZ,
			weight FLOAT,
			price BIGINT,
			packaging VARCHAR,
			additional_film VARCHAR,
			received_at TIMESTAMPTZ DEFAULT NULL,
			returned_at TIMESTAMPTZ DEFAULT NULL
		);
		CREATE TABLE IF NOT EXISTS orders_history (
			order_id BIGINT REFERENCES orders(order_id) PRIMARY KEY
		);
	`)
	require.NoError(b, err)

	// Очищаем и закрываем ресурсы после завершения тестов
	b.Cleanup(func() {
		pool.Exec(ctx, `DELETE FROM orders; DELETE FROM orders_history;`)
		cancel()
		pool.Close()
	})

	return repo, pool, ctx, cancel
}

func BenchmarkPgRepository_AddOrder(b *testing.B) {
	repo, _, ctx, _ := setupPostgresStorage(b)

	for i := 0; i < b.N; i++ {
		err := repo.AddOrder(ctx, int64(i+1), 100, time.Now(), time.Now().Add(24*time.Hour), 1.5, 500, "box", "yes")
		require.NoError(b, err)
	}
}

var orderID int64 = 1

func BenchmarkPgRepository_DeleteOrder(b *testing.B) {
	repo, _, ctx, _ := setupPostgresStorage(b)

	err := repo.AddOrder(ctx, orderID, 100, time.Now(), time.Now().Add(24*time.Hour), 1.5, 500, "box", "yes")
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		err = repo.DeleteOrder(ctx, orderID)
		require.NoError(b, err)
		// Добавляем обратно заказ для повторного удаления
		err = repo.AddOrder(ctx, orderID, 100, time.Now(), time.Now().Add(24*time.Hour), 1.5, 500, "box", "yes")
		require.NoError(b, err)
	}
}

func BenchmarkPgRepository_GiveOrders(b *testing.B) {
	repo, _, ctx, _ := setupPostgresStorage(b)

	err := repo.AddOrder(ctx, orderID, 100, time.Now(), time.Now().Add(24*time.Hour), 1.5, 500, "box", "yes")
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		err := repo.GiveOrders(ctx, []int64{orderID})
		require.NoError(b, err)
	}
}

func BenchmarkPgRepository_GetOrders(b *testing.B) {
	repo, _, ctx, _ := setupPostgresStorage(b)

	err := repo.AddOrder(ctx, orderID, 100, time.Now(), time.Now().Add(24*time.Hour), 1.5, 500, "box", "yes")
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		_, err := repo.GetOrders(ctx, 100)
		require.NoError(b, err)
	}
}

func BenchmarkPgRepository_AcceptReturn(b *testing.B) {
	repo, _, ctx, _ := setupPostgresStorage(b)

	err := repo.AddOrder(ctx, orderID, 100, time.Now(), time.Now().Add(24*time.Hour), 1.5, 500, "box", "yes")
	require.NoError(b, err)

	err = repo.GiveOrders(ctx, []int64{orderID})
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		err := repo.AcceptReturn(ctx, 100, orderID)
		require.NoError(b, err)
	}
}
