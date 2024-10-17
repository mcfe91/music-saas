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

	// Add the item to the cart
	if err := h.cartService.AddToCart(user, request.ProductID, request.Quantity); err != nil {
		http.Error(w, "Could not add item to cart", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item added to cart"})
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	var request struct {
		ProductID int `json:"product_id"`
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

	// Call the service to remove the item from the user's cart
	err := h.cartService.RemoveFromCart(user.ID, request.ProductID)
	if err != nil {
		if err.Error() == "cart not found" {
			http.Error(w, "Cart not found", http.StatusNotFound)
		} else if err.Error() == "item not found in cart" {
			http.Error(w, "Item not found in cart", http.StatusNotFound)
		} else {
			http.Error(w, "Could not remove item from cart", http.StatusInternalServerError)
		}
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Item removed from cart"})
}
