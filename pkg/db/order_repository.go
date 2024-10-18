package db

import (
	"context"
	"database/sql"
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

func (repo *PostgresOrderRepository) GetUserOrderItems(userID, orderID int) ([]*model.OrderItemWithProduct, error) {
	rows, err := repo.db.QueryContext(context.Background(),
		`SELECT oi.order_id, oi.product_id, oi.quantity, oi.price,
		        p.name AS product_name, p.description AS product_description
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = $1 AND o.user_id = $2`, orderID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.OrderItemWithProduct
	for rows.Next() {
		var item model.OrderItemWithProduct
		err := rows.Scan(&item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &item.ProductName, &item.ProductDescription)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	return items, nil
}
