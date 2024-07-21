-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products
(
    product_id   UUID        NOT NULL,
    name         TEXT        NOT NULL,
    description  TEXT        NOT NULL,
    available INTEGER     NOT NULL,
    quantity     INTEGER     NOT NULL,
    price        INTEGER     NOT NULL,
    location     VARCHAR(10) NOT NULL,
    PRIMARY KEY (product_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
