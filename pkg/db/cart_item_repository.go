package db

import (
	"context"
	"database/sql"
)

type PostgresCartItemRepository struct {
	db *sql.DB
}

func NewPostgresCartItemRepository(db *sql.DB) *PostgresCartItemRepository {
	return &PostgresCartItemRepository{db: db}
}

func (repo *PostgresCartItemRepository) AddCartItem(cartID, productID, quantity int) error {
	_, err := repo.db.ExecContext(context.Background(),
		"INSERT INTO cart_items (cart_id, product_id, quantity) VALUES ($1, $2, $3)",
		cartID, productID, quantity)
	return err
}

func (repo *PostgresCartItemRepository) RemoveCartItem(cartID, productID int) (bool, error) {
	result, err := repo.db.ExecContext(context.Background(),
		"DELETE FROM cart_items WHERE cart_id = $1 AND product_id = $2",
		cartID, productID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}
