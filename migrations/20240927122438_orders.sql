-- +goose Up
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
-- +goose Down
DROP TABLE IF EXISTS orders;