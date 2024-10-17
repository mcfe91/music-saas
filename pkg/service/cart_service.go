package service

import (
	"errors"
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

func (s *CartService) AddToCart(user *model.User, productID, quantity int) error {
	cart, err := s.GetCartByUserID(user.ID)
	// TODO: add specific error that cart doesn't exist
	if err != nil {
		cart, err = s.CreateCart(user.ID)
		if err != nil {
			return err
		}
	}
	return s.cartItemRepo.AddCartItem(cart.ID, productID, quantity)
}

func (s *CartService) GetCartByUserID(userID int) (*model.Cart, error) {
	return s.cartRepo.GetCartByUserID(userID)
}

func (s *CartService) RemoveFromCart(userID, productID int) error {
	cart, err := s.cartRepo.GetCartByUserID(userID)
	if err != nil {
		return err
	}
	if cart == nil {
		return errors.New("cart not found")
	}

	itemRemoved, err := s.cartItemRepo.RemoveCartItem(cart.ID, productID)
	if err != nil {
		return err
	}
	if !itemRemoved {
		return errors.New("item not found in cart")
	}

	return nil
}
