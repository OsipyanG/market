-- +goose Up
-- +goose StatementBegin
CREATE TABLE Products (
    product_id UUID NOT NULL,
    user_id UUID NOT NULL,
    quantity INTEGER NOT NULL,
    PRIMARY KEY (user_id, product_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS Products;
-- +goose StatementEnd
