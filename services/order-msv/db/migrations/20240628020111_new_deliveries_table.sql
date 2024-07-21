-- +goose Up
-- +goose StatementBegin
CREATE TYPE DELIVERY_STATUS AS ENUM ('assembly', 'pending', 'delivery', 'completed', 'cancelled');
CREATE TABLE IF NOT EXISTS deliveries
(
    order_id   UUID            NOT NULL,
    courier_id UUID            NOT NULL,
    status     DELIVERY_STATUS NOT NULL
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS deliveries;
DROP TYPE IF EXISTS DELIVERY_STATUS;
-- +goose StatementEnd
