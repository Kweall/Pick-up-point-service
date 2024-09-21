package json_file

import (
	"homework/storage/json_file"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setupStorage(b *testing.B) *json_file.Storage {
	file, err := os.CreateTemp("", "orders_bench_*.json")
	require.NoError(b, err)

	storage, err := json_file.NewStorage(file.Name())
	require.NoError(b, err)
	b.Cleanup(func() {
		os.Remove(file.Name())
	})
	return storage
}

func BenchmarkStorage_AddOrder(b *testing.B) {
	storage := setupStorage(b)

	order := &json_file.Order{
		ID:        1,
		ClientID:  100,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
		Weight:    1.5,
		Price:     500,
		Packaging: "Box",
	}

	for i := 0; i < b.N; i++ {
		err := storage.AddOrder(order)
		require.NoError(b, err)
	}
}

func BenchmarkStorage_DeleteOrderByID(b *testing.B) {
	storage := setupStorage(b)

	order := &json_file.Order{
		ID:        1,
		ClientID:  100,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
		Weight:    1.5,
		Price:     500,
		Packaging: "Box",
	}
	err := storage.AddOrder(order)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		err := storage.DeleteOrderByID(order.ID)
		require.NoError(b, err)

		err = storage.AddOrder(order)
		require.NoError(b, err)
	}
}

func BenchmarkStorage_GiveOrdersToClient(b *testing.B) {
	storage := setupStorage(b)

	order := &json_file.Order{
		ID:        1,
		ClientID:  100,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
		Weight:    1.5,
		Price:     500,
		Packaging: "Box",
	}
	err := storage.AddOrder(order)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		order.RecievedAt = time.Time{}
		err := storage.AddOrder(order)
		require.NoError(b, err)

		err = storage.GiveOrdersToClient([]int64{order.ID})
		require.NoError(b, err)
	}
}

func BenchmarkStorage_GetAll(b *testing.B) {
	storage := setupStorage(b)

	order := &json_file.Order{
		ID:        1,
		ClientID:  100,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(24 * time.Hour),
		Weight:    1.5,
		Price:     500,
		Packaging: "Box",
	}
	err := storage.AddOrder(order)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		_, err := storage.GetAll()
		require.NoError(b, err)
	}
}

func BenchmarkStorage_AcceptReturn(b *testing.B) {
	storage := setupStorage(b)

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
	err := storage.AddOrder(order)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		err := storage.AcceptReturn(order.ClientID, order.ID)
		require.NoError(b, err)

		order.ReturnedAt = time.Time{}
		err = storage.AddOrder(order)
		require.NoError(b, err)
	}
}
