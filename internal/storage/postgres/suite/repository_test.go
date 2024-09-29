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
		CREATE TABLE IF NOT EXISTS orders (
			order_id BIGINT PRIMARY KEY,
			client_id BIGINT,
			created_at TIMESTAMP,
			expired_at TIMESTAMP,
			received_at TIMESTAMP DEFAULT NULL,
			returned_at TIMESTAMP DEFAULT NULL,
			weight FLOAT,
			price INT,
			packaging VARCHAR,
			additional_film VARCHAR
		);
		CREATE TABLE IF NOT EXISTS orders_history (
			order_id BIGINT PRIMARY KEY
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

var orderID int64 = 10

func (suite *Suite) TestAddOrder() {
	err := suite.repo.AddOrder(suite.ctx, orderID, 20, time.Now(), time.Now().Add(24*time.Hour), 1.5, 500, "box", "yes")
	require.NoError(suite.T(), err)

	var count int
	err = suite.pool.QueryRow(suite.ctx, "SELECT COUNT(*) FROM orders WHERE order_id = $1", orderID).Scan(&count)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 1, count)
}

func (suite *Suite) TestDeleteOrder() {
	err := suite.repo.AddOrder(suite.ctx, orderID, 20, time.Now(), time.Now().Add(24*time.Hour), 1.5, 500, "box", "yes")
	require.NoError(suite.T(), err)

	err = suite.repo.DeleteOrder(suite.ctx, orderID)
	require.NoError(suite.T(), err)

	var count int
	err = suite.pool.QueryRow(suite.ctx, "SELECT COUNT(*) FROM orders WHERE order_id = $1", orderID).Scan(&count)
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), 0, count)
}

func (suite *Suite) TestGiveOrders() {
	err := suite.repo.AddOrder(suite.ctx, orderID, 20, time.Now(), time.Now().Add(24*time.Hour), 1.5, 500, "box", "yes")
	require.NoError(suite.T(), err)

	err = suite.repo.GiveOrders(suite.ctx, []int64{orderID})
	require.NoError(suite.T(), err)

	var receivedAt time.Time
	err = suite.pool.QueryRow(suite.ctx, "SELECT received_at FROM orders WHERE order_id = $1", orderID).Scan(&receivedAt)
	require.NoError(suite.T(), err)
	require.False(suite.T(), receivedAt.IsZero())
}

func (suite *Suite) TestAcceptReturn() {
	createdAt := time.Now()
	err := suite.repo.AddOrder(suite.ctx, orderID, 20, createdAt, createdAt.Add(24*time.Hour), 1.5, 500, "box", "yes")
	require.NoError(suite.T(), err)

	err = suite.repo.GiveOrders(suite.ctx, []int64{orderID})
	require.NoError(suite.T(), err)

	err = suite.repo.AcceptReturn(suite.ctx, 20, orderID)
	require.NoError(suite.T(), err)

	var returnedAt time.Time
	err = suite.pool.QueryRow(suite.ctx, "SELECT returned_at FROM orders WHERE order_id = $1", orderID).Scan(&returnedAt)
	require.NoError(suite.T(), err)
	require.False(suite.T(), returnedAt.IsZero())
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
