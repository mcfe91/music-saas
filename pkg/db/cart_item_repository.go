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
