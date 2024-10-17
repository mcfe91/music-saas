package api

import (
	"encoding/json"
	"music-saas/internal/middleware"
	"music-saas/pkg/model"
	"music-saas/pkg/service"
	"net/http"
)

type CartHandler struct {
	cartService *service.CartService
}

func NewCartHandler(cartService *service.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve user from context
	user, ok := r.Context().Value(middleware.ContextUserKey).(*model.User)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Retrieve the user's cart
	cart, err := h.cartService.GetCartByUserID(user.ID)
	if err != nil {
		// If the cart does not exist, create a new one
		cart, err = h.cartService.CreateCart(user.ID)
		if err != nil {
			http.Error(w, "Could not create cart", http.StatusInternalServerError)
			return
		}
	}

	// Add the item to the cart
	if err := h.cartService.AddToCart(cart.ID, request.ProductID, request.Quantity); err != nil {
		http.Error(w, "Could not add item to cart", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

// func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
// 	var request struct {
// 		ProductID int `json:"product_id"`
// 		Quantity  int `json:"quantity"`
// 	}

// 	// Decode the request body
// 	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 		http.Error(w, "Invalid request", http.StatusBadRequest)
// 		return
// 	}

// 	// Retrieve user from context
// 	user, ok := r.Context().Value(middleware.ContextUserKey).(*model.User)
// 	if !ok {
// 		http.Error(w, "User not found", http.StatusUnauthorized)
// 		return
// 	}

// 	// Create a cart if it doesn't exist (you might want to handle this differently)
// 	cart, err := h.cartService.CreateCart(user.ID)
// 	if err != nil {
// 		http.Error(w, "Could not create cart", http.StatusInternalServerError)
// 		return
// 	}

// 	// Add the item to the cart
// 	if err := h.cartService.AddToCart(cart.ID, request.ProductID, request.Quantity); err != nil {
// 		http.Error(w, "Could not add item to cart", http.StatusInternalServerError)
// 		return
// 	}

// 	// Respond with success
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(cart)
// }
