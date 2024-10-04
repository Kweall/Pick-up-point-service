-- +goose Up
CREATE TABLE IF NOT EXISTS orders_history (order_id BIGINT PRIMARY KEY NOT NULL);
-- +goose Down
DROP TABLE IF EXISTS orders_history;