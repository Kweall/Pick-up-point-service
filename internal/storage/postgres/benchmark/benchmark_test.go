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
		CREATE TYPE packaging AS ENUM ('box', 'bag', 'film');
		CREATE TABLE IF NOT EXISTS orders (
			order_id BIGINT PRIMARY KEY NOT NULL,
			client_id BIGINT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			expired_at TIMESTAMP NOT NULL,
			received_at TIMESTAMP DEFAULT NULL,
			returned_at TIMESTAMP DEFAULT NULL,
			weight FLOAT NOT NULL,
			price BIGINT NOT NULL,
			packaging packaging NOT NULL,
			additional_film BOOLEAN NOT NULL
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
	createdAt := time.Now().Truncate(time.Minute)
	parsedDate := createdAt.Add(24 * time.Hour)
	req := &postgres.Order{
		ClientID:       1,
		OrderID:        100,
		CreatedAt:      &createdAt,
		ExpiredAt:      &parsedDate,
		Weight:         1.5,
		Price:          500,
		Packaging:      "box",
		AdditionalFilm: true,
	}
	for i := 0; i < b.N; i++ {
		err := repo.AddOrder(ctx, req)
		require.NoError(b, err)
	}
}

func BenchmarkPgRepository_DeleteOrder(b *testing.B) {
	repo, _, ctx, _ := setupPostgresStorage(b)
	createdAt := time.Now().Truncate(time.Minute)
	parsedDate := createdAt.Add(24 * time.Hour)
	req := &postgres.Order{
		ClientID:       1,
		OrderID:        100,
		CreatedAt:      &createdAt,
		ExpiredAt:      &parsedDate,
		Weight:         1.5,
		Price:          500,
		Packaging:      "box",
		AdditionalFilm: true,
	}
	err := repo.AddOrder(ctx, req)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		err = repo.DeleteOrder(ctx, req.OrderID)
		require.NoError(b, err)
		// Добавляем обратно заказ для повторного удаления
		err = repo.AddOrder(ctx, req)
		require.NoError(b, err)
	}
}

func BenchmarkPgRepository_GiveOrders(b *testing.B) {
	repo, _, ctx, _ := setupPostgresStorage(b)

	createdAt := time.Now().Truncate(time.Minute)
	parsedDate := createdAt.Add(24 * time.Hour)
	req := &postgres.Order{
		ClientID:       1,
		OrderID:        100,
		CreatedAt:      &createdAt,
		ExpiredAt:      &parsedDate,
		Weight:         1.5,
		Price:          500,
		Packaging:      "box",
		AdditionalFilm: true,
	}
	err := repo.AddOrder(ctx, req)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		err := repo.GiveOrders(ctx, []int64{req.OrderID})
		require.NoError(b, err)
	}
}

func BenchmarkPgRepository_GetOrders(b *testing.B) {
	repo, _, ctx, _ := setupPostgresStorage(b)

	createdAt := time.Now().Truncate(time.Minute)
	parsedDate := createdAt.Add(24 * time.Hour)
	req := &postgres.Order{
		ClientID:       1,
		OrderID:        100,
		CreatedAt:      &createdAt,
		ExpiredAt:      &parsedDate,
		Weight:         1.5,
		Price:          500,
		Packaging:      "box",
		AdditionalFilm: true,
	}
	err := repo.AddOrder(ctx, req)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		_, err := repo.GetOrders(ctx, req.ClientID)
		require.NoError(b, err)
	}
}

func BenchmarkPgRepository_AcceptReturn(b *testing.B) {
	repo, _, ctx, _ := setupPostgresStorage(b)

	createdAt := time.Now().Truncate(time.Minute)
	parsedDate := createdAt.Add(24 * time.Hour)
	req := &postgres.Order{
		ClientID:       1,
		OrderID:        100,
		CreatedAt:      &createdAt,
		ExpiredAt:      &parsedDate,
		Weight:         1.5,
		Price:          500,
		Packaging:      "box",
		AdditionalFilm: true,
	}
	err := repo.AddOrder(ctx, req)
	require.NoError(b, err)

	err = repo.GiveOrders(ctx, []int64{req.OrderID})
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		err := repo.AcceptReturn(ctx, req.ClientID, req.OrderID)
		require.NoError(b, err)
	}
}
