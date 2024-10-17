package service

import (
	"music-saas/pkg/db"
	"music-saas/pkg/model"
)

type CartService struct {
	cartRepo     *db.PostgresCartRepository
	cartItemRepo *db.PostgresCartItemRepository
}

func NewCartService(cartRepo *db.PostgresCartRepository, cartItemRepo *db.PostgresCartItemRepository) *CartService {
	return &CartService{cartRepo: cartRepo, cartItemRepo: cartItemRepo}
}

func (s *CartService) CreateCart(userID int) (*model.Cart, error) {
	return s.cartRepo.CreateCart(userID)
}

func (s *CartService) AddToCart(cartID, productID, quantity int) error {
	return s.cartItemRepo.AddCartItem(cartID, productID, quantity)
}

func (s *CartService) GetCartByUserID(userID int) (*model.Cart, error) {
	return s.cartRepo.GetCartByUserID(userID)
}
