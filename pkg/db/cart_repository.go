package db

import (
	"context"
	"database/sql"
	"music-saas/pkg/model"
)

type PostgresCartRepository struct {
	db *sql.DB
}

func NewPostgresCartRepository(db *sql.DB) *PostgresCartRepository {
	return &PostgresCartRepository{db: db}
}

func (repo *PostgresCartRepository) CreateCart(userID int) (*model.Cart, error) {
	cart := &model.Cart{UserID: userID}
	err := repo.db.QueryRowContext(context.Background(),
		"INSERT INTO carts (user_id) VALUES ($1) RETURNING id, created_at, updated_at",
		userID).Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (repo *PostgresCartRepository) GetCartByUserID(userID int) (*model.Cart, error) {
	cart := &model.Cart{}
	row := repo.db.QueryRowContext(context.Background(), "SELECT id, user_id FROM carts WHERE user_id = $1", userID)
	err := row.Scan(&cart.ID, &cart.UserID)
	if err != nil {
		return nil, err
	}
	return cart, nil
}
