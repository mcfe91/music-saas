package model

type OrderItem struct {
	ID        int     `json:"id"`
	OrderID   int     `json:"order_id"`   // Foreign key to Order
	ProductID int     `json:"product_id"` // Foreign key to Product
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"` // Price per unit at the time of order
}

type OrderItemWithProduct struct {
	OrderID            int     `json:"order_id"`
	ProductID          int     `json:"product_id"`
	Quantity           int     `json:"quantity"`
	Price              float64 `json:"price"`
	ProductName        string  `json:"product_name"`
	ProductDescription string  `json:"product_description"`
}
