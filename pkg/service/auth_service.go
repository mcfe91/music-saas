package service

import (
	"errors"
	"music-saas/pkg/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Signup(username, password string) (*model.User, error) {
	user := &model.User{Username: username, Password: password}
	err := user.HashPassword()
	if err != nil {
		return nil, err
	}
	// Store in database here
	return user, nil
}

func Login(username, password string) (string, error) {
	// Fetch user in database here
	user := &model.User{Username: username, Password: "$2a$10$vPYge9uhVepnUiuS7LIEMufXs9ukhb56im70K9oNRYKBLxZN7zbyW"}

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
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
