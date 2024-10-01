-- +goose Up
CREATE TYPE packaging AS ENUM ('box', 'bag', 'film');
CREATE TABLE IF NOT EXISTS orders (
    order_id BIGINT PRIMARY KEY NOT NULL,
    client_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    received_at TIMESTAMP DEFAULT NULL,
    returned_at TIMESTAMP DEFAULT NULL,
    weight FLOAT NOT NULL,
    price INT NOT NULL,
    packaging packaging NOT NULL,
    additional_film BOOLEAN NOT NULL
);
-- +goose Down
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS packaging;