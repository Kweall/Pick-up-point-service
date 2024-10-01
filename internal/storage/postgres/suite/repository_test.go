package postgres

import (
	"context"
	"testing"
	"time"

	"homework/internal/storage/postgres"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	repo      *postgres.PgRepository
	pool      *pgxpool.Pool
	txManager *postgres.TxManager
	ctx       context.Context
	cancel    context.CancelFunc
}

// Перед каждым тестом создаем новое соединение с базой данных
func (suite *Suite) SetupTest() {
	var err error
	suite.ctx, suite.cancel = context.WithCancel(context.Background())

	suite.pool, err = pgxpool.Connect(suite.ctx, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	require.NoError(suite.T(), err)

	suite.txManager = postgres.NewTxManager(suite.pool)

	suite.repo = postgres.NewPgRepository(suite.txManager)

	_, err = suite.pool.Exec(suite.ctx, `
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
	require.NoError(suite.T(), err)
}

// После каждого теста удаляем временное соединение
func (suite *Suite) TearDownTest() {
	_, err := suite.pool.Exec(suite.ctx, `
		DROP TABLE IF EXISTS orders_history;
		DROP TABLE IF EXISTS orders;
	`)
	require.NoError(suite.T(), err)

	suite.cancel()
	suite.pool.Close()
}

func (suite *Suite) TestAddOrder() {
	createdAt := time.Now().Truncate(time.Minute)
	parsedDate := createdAt.Add(24 * time.Hour)
	req := &postgres.Order{
		ClientID:       20,
		OrderID:        10,
		CreatedAt:      &createdAt,
		ExpiredAt:      &parsedDate,
		Weight:         1.5,
		Price:          500,
		Packaging:      "box",
		AdditionalFilm: true,
	}
	err := suite.repo.AddOrder(suite.ctx, req)
	require.NoError(suite.T(), err)

	var count int
	err = suite.pool.QueryRow(suite.ctx, "SELECT COUNT(*) FROM orders WHERE order_id = $1", req.OrderID).Scan(&count)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 1, count)
}

func (suite *Suite) TestDeleteOrder() {
	createdAt := time.Now().Truncate(time.Minute)
	parsedDate := createdAt.Add(24 * time.Hour)
	req := &postgres.Order{
		ClientID:       20,
		OrderID:        10,
		CreatedAt:      &createdAt,
		ExpiredAt:      &parsedDate,
		Weight:         1.5,
		Price:          500,
		Packaging:      "box",
		AdditionalFilm: true,
	}
	err := suite.repo.AddOrder(suite.ctx, req)
	require.NoError(suite.T(), err)

	err = suite.repo.DeleteOrder(suite.ctx, req.OrderID)
	require.NoError(suite.T(), err)

	var count int
	err = suite.pool.QueryRow(suite.ctx, "SELECT COUNT(*) FROM orders WHERE order_id = $1", req.OrderID).Scan(&count)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 0, count)
}

func (suite *Suite) TestGiveOrders() {
	createdAt := time.Now().Truncate(time.Minute)
	parsedDate := createdAt.Add(24 * time.Hour)
	req := &postgres.Order{
		ClientID:       20,
		OrderID:        10,
		CreatedAt:      &createdAt,
		ExpiredAt:      &parsedDate,
		Weight:         1.5,
		Price:          500,
		Packaging:      "box",
		AdditionalFilm: true,
	}
	err := suite.repo.AddOrder(suite.ctx, req)
	require.NoError(suite.T(), err)

	err = suite.repo.GiveOrders(suite.ctx, []int64{req.OrderID})
	require.NoError(suite.T(), err)

	var receivedAt time.Time
	err = suite.pool.QueryRow(suite.ctx, "SELECT received_at FROM orders WHERE order_id = $1", req.OrderID).Scan(&receivedAt)
	require.NoError(suite.T(), err)
	require.False(suite.T(), receivedAt.IsZero())
}

func (suite *Suite) TestAcceptReturn() {
	createdAt := time.Now().Truncate(time.Minute)
	parsedDate := createdAt.Add(24 * time.Hour)
	req := &postgres.Order{
		ClientID:       20,
		OrderID:        10,
		CreatedAt:      &createdAt,
		ExpiredAt:      &parsedDate,
		Weight:         1.5,
		Price:          500,
		Packaging:      "box",
		AdditionalFilm: true,
	}
	err := suite.repo.AddOrder(suite.ctx, req)
	require.NoError(suite.T(), err)

	err = suite.repo.GiveOrders(suite.ctx, []int64{req.OrderID})
	require.NoError(suite.T(), err)

	err = suite.repo.AcceptReturn(suite.ctx, 20, req.OrderID)
	require.NoError(suite.T(), err)

	var returnedAt time.Time
	err = suite.pool.QueryRow(suite.ctx, "SELECT returned_at FROM orders WHERE order_id = $1", req.OrderID).Scan(&returnedAt)
	require.NoError(suite.T(), err)
	require.False(suite.T(), returnedAt.IsZero())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
