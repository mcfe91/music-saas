package api

import (
	"encoding/json"
	"log"
	"music-saas/internal/middleware"
	"music-saas/pkg/model"
	"music-saas/pkg/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type OrderHandler struct {
	orderService *service.OrderService
}

func NewOrderHandler(orderService *service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order model.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.orderService.CreateOrder(&order)
	if err != nil {
		log.Println("Error creating order:", err)
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *OrderHandler) GetOrderItems(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(middleware.ContextUserKey).(*model.User)
	if !ok {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	orderItems, err := h.orderService.GetOrderItems(orderID, user.ID)
	if err != nil {
		if err.Error() == "unauthorized access to order items" {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			http.Error(w, "Failed to retrieve order items", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(orderItems)
}
