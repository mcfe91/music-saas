package service

import (
	"errors"
	"music-saas/pkg/db"
	"music-saas/pkg/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	ID       int    `json:"id"`
	jwt.StandardClaims
}

type AuthService struct {
	userRepo db.UserRepository
	jwtKey   []byte
}

func NewAuthService(repo db.UserRepository, jwtKey []byte) *AuthService {
	return &AuthService{userRepo: repo, jwtKey: jwtKey}
}

func (s *AuthService) Signup(username, password string) (*model.User, error) {
	user := &model.User{Username: username, Password: password}
	err := user.HashPassword()
	if err != nil {
		return nil, err
	}
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	// TODO: User doesn't return correct id from db, fix or just remove...
	return user, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if !user.CheckPassword(password) {
		return "", errors.New("invalid credentials")
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		Username: user.Username,
		ID:       user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) GetUserByUsername(username string) (*model.User, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return user, nil
}
