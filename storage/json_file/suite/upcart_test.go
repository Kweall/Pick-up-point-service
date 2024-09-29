package json_file

import (
	"homework/storage/json_file"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Определяем структуру для тестового сьюта
type Suite struct {
	suite.Suite
	storage *json_file.Storage
	file    *os.File
}

// Перед каждым тестом создаем временное хранилище
func (suite *Suite) SetupTest() {
	var err error
	suite.file, err = os.CreateTemp("", "orders_*.json")
	require.NoError(suite.T(), err)

	suite.storage, err = json_file.NewStorage(suite.file.Name())
	require.NoError(suite.T(), err)
}

// После каждого теста удаляем временный файл
func (suite *Suite) TearDownTest() {
	os.Remove(suite.file.Name())
}

// Тест на добавление заказа
func (suite *Suite) TestAddOrder() {
	order := &json_file.Order{
		ID:        1,
		ClientID:  100,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
		Weight:    1.5,
		Price:     500,
		Packaging: "Box",
	}

	err := suite.storage.AddOrder(order)
	require.NoError(suite.T(), err)

	orders, err := suite.storage.GetAll()
	require.NoError(suite.T(), err)
	require.Equal(suite.T(), order, orders[1])
}

// Тест на удаление заказа по ID
func (suite *Suite) TestDeleteOrderByID() {
	order := &json_file.Order{
		ID:        1,
		ClientID:  100,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
		Weight:    1.5,
		Price:     500,
		Packaging: "Box",
	}

	err := suite.storage.AddOrder(order)
	require.NoError(suite.T(), err)

	err = suite.storage.DeleteOrderByID(1)
	require.NoError(suite.T(), err)

	orders, err := suite.storage.GetAll()
	require.NoError(suite.T(), err)
	require.Empty(suite.T(), orders)
}

// Тест на выдачу заказов клиенту
func (suite *Suite) TestGiveOrdersToClient() {
	order := &json_file.Order{
		ID:        1,
		ClientID:  100,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
		Weight:    1.5,
		Price:     500,
		Packaging: "Box",
	}

	err := suite.storage.AddOrder(order)
	require.NoError(suite.T(), err)

	err = suite.storage.GiveOrdersToClient([]int64{1})
	require.NoError(suite.T(), err)

	orders, err := suite.storage.GetAll()
	require.NoError(suite.T(), err)
	require.False(suite.T(), orders[1].RecievedAt.IsZero())
}

// Тест на возврат заказа
func (suite *Suite) TestAcceptReturn() {
	order := &json_file.Order{
		ID:         1,
		ClientID:   100,
		CreatedAt:  time.Now(),
		ExpiredAt:  time.Now().Add(24 * time.Hour),
		RecievedAt: time.Now(),
		Weight:     1.5,
		Price:      500,
		Packaging:  "Box",
	}

	err := suite.storage.AddOrder(order)
	require.NoError(suite.T(), err)

	err = suite.storage.AcceptReturn(100, order.ID)
	require.NoError(suite.T(), err)

	orders, err := suite.storage.GetAll()
	require.NoError(suite.T(), err)
	require.False(suite.T(), orders[1].ReturnedAt.IsZero())
}

// Запускаем тесты
func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
