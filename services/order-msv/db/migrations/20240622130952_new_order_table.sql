-- +goose Up
-- +goose StatementBegin
CREATE TYPE ORDER_STATUS AS ENUM ('assembly', 'pending', 'delivery', 'completed', 'cancelled');
CREATE TABLE IF NOT EXISTS orders
(
    order_id    UUID PRIMARY KEY,
    customer_id UUID         NOT NULL,
    status      ORDER_STATUS NOT NULL DEFAULT 'assembly',
    address     TEXT         NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TYPE IF EXISTS ORDER_STATUS;
-- +goose StatementEnd
