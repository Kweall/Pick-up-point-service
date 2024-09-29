package postgres

import (
	"time"
)

type Order struct {
	OrderID        int64      `db:"order_id"`
	ClientID       int64      `db:"client_id"`
	CreatedAt      *time.Time `db:"created_at"`
	ExpiredAt      *time.Time `db:"expired_at"`
	ReceivedAt     *time.Time `db:"received_at"`
	ReturnedAt     *time.Time `db:"returned_at"`
	Weight         float64    `db:"weight"`
	Price          int64      `db:"price"`
	Packaging      string     `db:"packaging"`
	AdditionalFilm string     `db:"additional_film"`
}

type OrderHistory struct {
	OrderID int64 `db:"order_id"`
}
