package db

import (
	"context"
	"database/sql"
	"music-saas/pkg/model"
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

func (repo *PostgresCartItemRepository) GetCartItemByProductID(cartID, productID int) (*model.CartItem, error) {
	item := &model.CartItem{}
	row := repo.db.QueryRowContext(context.Background(), "SELECT id, cart_id, product_id, quantity FROM cart_items WHERE cart_id = $1 AND product_id = $2", cartID, productID)
	err := row.Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No item found
		}
		return nil, err
	}
	return item, nil
}

func (repo *PostgresCartItemRepository) UpdateCartItemQuantity(id, quantity int) error {
	_, err := repo.db.ExecContext(context.Background(),
		"UPDATE cart_items SET quantity = $1 WHERE id = $2",
		quantity, id)
	return err
}
