-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE PROCEDURE add_product(user_id UUID, product_id UUID, quantity INT) AS $$
DECLARE
    product_count INT;
BEGIN
    SELECT COUNT(*) INTO product_count FROM Products AS p
    WHERE p.user_id = add_product.user_id AND p.product_id = add_product.product_id;

    IF product_count > 0 THEN
        UPDATE Products AS p SET quantity = p.quantity + add_product.quantity
        WHERE p.user_id = add_product.user_id AND p.product_id = add_product.product_id;
    ELSE
        INSERT INTO Products (user_id, product_id, quantity)
        VALUES (add_product.user_id, add_product.product_id, add_product.quantity);
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION delete_product(user_id UUID, product_id UUID, quantity INT) RETURNS INT AS $$
DECLARE
    current_quantity INT;
    rows_affected INT;
BEGIN
    SELECT p.quantity INTO current_quantity FROM Products AS p
    WHERE p.user_id = delete_product.user_id AND p.product_id = delete_product.product_id;

    IF current_quantity <= quantity THEN
        DELETE FROM Products AS p WHERE p.user_id = delete_product.user_id AND p.product_id = delete_product.product_id;
        GET DIAGNOSTICS rows_affected = ROW_COUNT;
    ELSE
        UPDATE Products AS p SET quantity = p.quantity - delete_product.quantity
        WHERE p.user_id = delete_product.user_id AND p.product_id = delete_product.product_id;
        GET DIAGNOSTICS rows_affected = ROW_COUNT;
    END IF;

    RETURN rows_affected;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP PROCEDURE IF EXISTS add_product(UUID, UUID, INT);
DROP FUNCTION IF EXISTS delete_product(UUID, UUID, INT);

-- +goose StatementEnd
