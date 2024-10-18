package db

import (
	"context"
	"database/sql"
	"errors"
	"music-saas/pkg/model"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (repo *PostgresOrderRepository) CreateOrder(order *model.Order) error {
	// Start a transaction
	tx, err := repo.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}

	// Insert the order
	err = tx.QueryRowContext(
		context.Background(),
		"INSERT INTO orders (user_id, status) VALUES ($1, $2) RETURNING id",
		order.UserID, order.Status,
	).Scan(&order.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Insert order items
	for _, item := range order.Items {
		_, err := tx.ExecContext(
			context.Background(),
			"INSERT INTO order_items (order_id, product_id, quantity, price) VALUES ($1, $2, $3, $4)",
			order.ID, item.ProductID, item.Quantity, item.Price,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

func (repo *PostgresOrderRepository) GetOrderByID(id int) (*model.Order, error) {
	order := &model.Order{}
	row := repo.db.QueryRowContext(context.Background(),
		"SELECT id, user_id, created_at FROM orders WHERE id = $1", id)
	err := row.Scan(&order.ID, &order.UserID, &order.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	order.Items, err = repo.GetOrderItems(order.ID)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (repo *PostgresOrderRepository) GetOrderItems(orderID int) ([]*model.OrderItem, error) {
	rows, err := repo.db.QueryContext(context.Background(),
		"SELECT id, order_id, product_id, quantity, price FROM order_items WHERE order_id = $1", orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.OrderItem
	for rows.Next() {
		var item model.OrderItem
		err := rows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}

func (repo *PostgresOrderRepository) GetOrdersByUserID(userID int) ([]*model.Order, error) {
	rows, err := repo.db.QueryContext(context.Background(),
		"SELECT id, user_id, created_at FROM orders WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		// Retrieve order items
		order.Items, err = repo.GetOrderItems(order.ID)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	return orders, nil
}

func (repo *PostgresOrderRepository) GetOrderByOrderID(orderID int) (*model.Order, error) {
	order := &model.Order{}

	query := `SELECT id, user_id, status, created_at FROM orders WHERE id = $1`
	row := repo.db.QueryRow(query, orderID)

	err := row.Scan(&order.ID, &order.UserID, &order.Status, &order.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	return order, nil
}
