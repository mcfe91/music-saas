package db

import (
	"context"
	"database/sql"
	"errors"
	"music-saas/pkg/model"
)

// TODO: rename symbols to match others
type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (repo *PostgresUserRepository) CreateUser(user *model.User) error {
	_, err := repo.db.ExecContext(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresUserRepository) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	row := repo.db.QueryRowContext(context.Background(), "SELECT username, password, is_admin FROM users WHERE username = $1", username)
	err := row.Scan(&user.Username, &user.Password, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}
