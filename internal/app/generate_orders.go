package app

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GenerateFakeOrders(db *pgxpool.Pool, numOrders int) error {
	for i := 0; i < numOrders; i++ {
		orderID := int64(i + 1)
		clientID := gofakeit.Number(1000, 999999)
		createdAt := gofakeit.DateRange(time.Now().Add(-24*time.Hour), time.Now())
		expiredAt := createdAt.Add(14 * 24 * time.Hour)

		var weight float64
		var packaging string
		var additionalFilm string

		for {
			packaging = gofakeit.RandomString([]string{"box", "bag", "film"})
			weight = math.Round(gofakeit.Float64Range(0.5, 20.0)*1000) / 1000

			if (packaging == "bag" && weight < 10) || (packaging == "box" && weight < 20) || packaging == "film" {
				if packaging == "film" {
					additionalFilm = "no"
				} else {
					additionalFilm = gofakeit.RandomString([]string{"yes", "no"})
				}
				break
			}
		}

		_, err := db.Exec(context.Background(), `
            INSERT INTO orders (order_id, client_id, created_at, expired_at, weight, price, packaging, additional_film)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        `, orderID, clientID, createdAt, expiredAt, weight, gofakeit.Number(100, 10000), packaging, additionalFilm)
		if err != nil {
			return fmt.Errorf("failed to insert order: %w", err)
		}
	}
	return nil
}

func ClearTables(db *pgxpool.Pool) error {
	query := `
        DELETE FROM orders;
        DELETE FROM orders_history;
    `
	_, err := db.Exec(context.Background(), query)
	if err != nil {
		return fmt.Errorf("failed to clear tables: %w", err)
	}
	return nil
}
