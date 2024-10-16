package service

import (
	"errors"
	"music-saas/pkg/db"
	"music-saas/pkg/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TODO: move to env file, then secure server.
var JwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthService struct {
	userRepo db.UserRepository
}

func NewAuthService(repo db.UserRepository) *AuthService {
	return &AuthService{userRepo: repo}
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

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
