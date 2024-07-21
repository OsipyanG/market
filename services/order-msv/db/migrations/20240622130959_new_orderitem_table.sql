-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orderitems
(
    order_id   UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity   INT  NOT NULL,
    price      INT  NOT NULL,
    PRIMARY KEY (order_id, product_id),
    FOREIGN KEY (order_id) REFERENCES orders (order_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orderitems;
-- +goose StatementEnd
