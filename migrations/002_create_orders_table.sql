-- +goose Up
CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    number VARCHAR(255) NOT NULL UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users(id),
    status order_status NOT NULL DEFAULT 'NEW',
    accrual BIGINT, -- сумма начисленных баллов в копейках
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_number ON orders(number);

-- +goose Down
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS order_status; 