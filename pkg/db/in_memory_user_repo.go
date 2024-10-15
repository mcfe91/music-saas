package db

import (
	"errors"
	"music-saas/pkg/model"
)

type InMemoryUserRepository struct {
	users map[string]*model.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*model.User),
	}
}

func (repo *InMemoryUserRepository) CreateUser(user *model.User) error {
	if _, exists := repo.users[user.Username]; exists {
		return errors.New("user already exists")
	}
	repo.users[user.Username] = user
	return nil
}

func (repo *InMemoryUserRepository) GetUserByUsername(username string) (*model.User, error) {
	user, exists := repo.users[username]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
