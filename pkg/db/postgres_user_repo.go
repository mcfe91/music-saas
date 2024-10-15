package db

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"music-saas/pkg/model"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (repo *PostgresUserRepository) CreateUser(user *model.User) error {
	_, err := repo.db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresUserRepository) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	row := repo.db.QueryRow("SELECT username, password FROM users WHERE username = $1", username)
	err := row.Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}
