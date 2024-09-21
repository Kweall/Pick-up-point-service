package json_file

import (
	"homework/storage/json_file"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestE2EAcceptOrderReturn(t *testing.T) {

	file, err := os.CreateTemp("", "orders_e2e_*.json")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	storage, err := json_file.NewStorage(file.Name())
	require.NoError(t, err)

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

	err = storage.AddOrder(order)
	require.NoError(t, err)

	err = storage.AcceptReturn(order.ClientID, order.ID)
	require.NoError(t, err)

	orders, err := storage.GetAll()
	require.NoError(t, err)
	require.False(t, orders[order.ID].ReturnedAt.IsZero())
}
