-- +goose Up
CREATE TABLE IF NOT EXISTS orders_history (order_id BIGINT PRIMARY KEY);
-- +goose Down
DROP TABLE IF EXISTS orders_history;