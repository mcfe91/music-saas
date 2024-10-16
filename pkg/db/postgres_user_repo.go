package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"music-saas/pkg/model"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (repo *PostgresUserRepository) CreateUser(user *model.User) error {
	_, err := repo.db.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresUserRepository) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	row := repo.db.QueryRow(context.Background(), "SELECT username, password FROM users WHERE username = $1", username)
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}
