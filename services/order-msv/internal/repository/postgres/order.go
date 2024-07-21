package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/OsipyanG/market/services/order-msv/internal/model"
	"github.com/OsipyanG/market/services/order-msv/internal/repository"
	"github.com/OsipyanG/market/services/order-msv/pkg/errwrap"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: db}
}

func (or *OrderRepository) CreateOrder(ctx context.Context, order model.Order) error {
	transaction, err := or.pool.Begin(ctx)
	if err != nil {
		return errwrap.Wrap(repository.ErrCreateOrder, err)
	}

	defer func() {
		err := transaction.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) && !errors.Is(err, pgx.ErrTxCommitRollback) {
			slog.Error("Error rolling back transaction", "error:", err)
		}
	}()

	query := `
		INSERT INTO orders (order_id, customer_id, address)
		VALUES ($1, $2, $3)
	`

	_, err = transaction.Exec(ctx, query, order.ID, order.CustomerID, order.Address)
	if err != nil {
		return errwrap.Wrap(repository.ErrCreateOrder, err)
	}

	for _, item := range order.Items {
		query := `
			INSERT INTO orderitems (order_id, product_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`

		_, err = transaction.Exec(ctx, query, order.ID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return errwrap.Wrap(repository.ErrCreateOrder, err)
		}
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return errwrap.Wrap(repository.ErrCreateOrder, err)
	}

	return nil
}

func (or *OrderRepository) GetOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	query := `
		SELECT o.order_id,
			o.customer_id,
			o.status,
			address,
			o.created_at,
			o.updated_at,
			i.product_id,
			i.quantity,
			i.price
		FROM orders o
				JOIN orderitems i ON o.order_id = i.order_id
		WHERE o.order_id = $1
    `

	rows, err := or.pool.Query(ctx, query, orderID)
	if err != nil {
		return nil, errwrap.Wrap(repository.ErrGetOrder, err)
	}
	defer rows.Close()

	order := &model.Order{}

	for rows.Next() {
		item := model.OrderItem{}

		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.Status,
			&order.Address,
			&order.CreatedAt,
			&order.UpdatedAt,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, errwrap.Wrap(repository.ErrGetOrder, err)
		}

		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, errwrap.Wrap(repository.ErrGetOrder, err)
	}

	return order, nil
}

func (or *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error {
	query := `
		UPDATE orders
		SET status     = $1,
			updated_at = CURRENT_TIMESTAMP
		WHERE order_id = $2
	`

	cmdTag, err := or.pool.Exec(ctx, query, status, orderID)
	if err != nil {
		return errwrap.Wrap(repository.ErrUpdateOrderStatus, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errwrap.Wrap(repository.ErrUpdateOrderStatus, repository.ErrNoRowsAffected)
	}

	return nil
}

func (or *OrderRepository) GetAllOrders(ctx context.Context, customerID uuid.UUID) ([]model.Order, error) {
	query := `
		SELECT o.order_id,
			o.customer_id,
			o.status,
			o.address,
			o.created_at,
			o.updated_at,
			i.product_id,
			i.quantity,
			i.price
		FROM orders o
				LEFT JOIN orderitems i ON o.order_id = i.order_id
		WHERE o.customer_id = $1
	`

	rows, err := or.pool.Query(ctx, query, customerID)
	if err != nil {
		return nil, errwrap.Wrap(repository.ErrGetAllOrders, err)
	}
	defer rows.Close()

	orderMap := make(map[uuid.UUID]*model.Order)

	for rows.Next() {
		order := &model.Order{}
		item := &model.OrderItem{}

		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.Status,
			&order.Address,
			&order.CreatedAt,
			&order.UpdatedAt,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, errwrap.Wrap(repository.ErrGetAllOrders, err)
		}

		if _, ok := orderMap[order.ID]; !ok {
			order.Items = append(order.Items, *item)
			orderMap[order.ID] = order
		} else {
			order := orderMap[order.ID]
			order.Items = append(order.Items, *item)
			orderMap[order.ID] = order
		}
	}

	orders := make([]model.Order, 0, len(orderMap))
	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	return orders, nil
}

func (or *OrderRepository) GetAllDeliveries(ctx context.Context) ([]model.Delivery, error) {
	query := `
		SELECT order_id, courier_id
		FROM deliveries
	`

	rows, err := or.pool.Query(ctx, query)
	if err != nil {
		return nil, errwrap.Wrap(repository.ErrGetAllDeliveries, err)
	}
	defer rows.Close()

	deliveries := make([]model.Delivery, 0)

	for rows.Next() {
		delivery := model.Delivery{}

		err := rows.Scan(&delivery.OrderID, &delivery.CourierID)
		if err != nil {
			return nil, errwrap.Wrap(repository.ErrGetAllDeliveries, err)
		}

		deliveries = append(deliveries, delivery)
	}

	return deliveries, nil
}

func (or *OrderRepository) GetAllPendingOrders(ctx context.Context) ([]model.Order, error) {
	query := `
		SELECT o.order_id,
			o.customer_id,
			o.status,
			o.address,
			o.created_at,
			o.updated_at,
			i.product_id,
			i.quantity,
			i.price
		FROM orders o
				LEFT JOIN orderitems i ON o.order_id = i.order_id
		WHERE o.status = 'pending'
	`

	rows, err := or.pool.Query(ctx, query)
	if err != nil {
		return nil, errwrap.Wrap(repository.ErrGetAllPendingOrders, err)
	}
	defer rows.Close()

	orderMap := make(map[uuid.UUID]*model.Order)

	for rows.Next() {
		order := &model.Order{}
		item := &model.OrderItem{}

		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.Status,
			&order.Address,
			&order.CreatedAt,
			&order.UpdatedAt,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, errwrap.Wrap(repository.ErrGetAllPendingOrders, err)
		}

		if _, ok := orderMap[order.ID]; !ok {
			order.Items = append(order.Items, *item)
			orderMap[order.ID] = order
		} else {
			order := orderMap[order.ID]
			order.Items = append(order.Items, *item)
			orderMap[order.ID] = order
		}
	}

	orders := make([]model.Order, 0, len(orderMap))
	for _, order := range orderMap {
		orders = append(orders, *order)
	}

	return orders, nil
}

func (or *OrderRepository) CreateDelivery(ctx context.Context, delivery model.Delivery) error {
	query := `
		INSERT INTO deliveries (order_id, courier_id)
		VALUES ($1, $2)
	`

	_, err := or.pool.Exec(ctx, query, delivery.OrderID, delivery.CourierID)
	if err != nil {
		return errwrap.Wrap(repository.ErrCreateDelivery, err)
	}

	return nil
}

func (or *OrderRepository) UpdateDeliveryStatus(ctx context.Context, orderID uuid.UUID, status string) error {
	query := `
		UPDATE deliveries
		SET status     = $1,
		where order_id = $2
	`

	cmdTag, err := or.pool.Exec(ctx, query, status, orderID)
	if err != nil {
		return errwrap.Wrap(repository.ErrUpdateDeliveryStatus, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errwrap.Wrap(repository.ErrUpdateDeliveryStatus, repository.ErrNoRowsAffected)
	}

	return nil
}
