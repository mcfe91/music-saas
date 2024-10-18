package service

import (
	"fmt"
	"music-saas/pkg/db"
	"music-saas/pkg/model"
)

type OrderService struct {
	orderRepo   *db.PostgresOrderRepository
	productRepo *db.PostgresProductRepository
}

func NewOrderService(orderRepo *db.PostgresOrderRepository, productRepo *db.PostgresProductRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo, productRepo: productRepo}
}

func (s *OrderService) CreateOrder(order *model.Order) error {
	if err := s.validateProducts(order.Items); err != nil {
		return err
	}

	return s.orderRepo.CreateOrder(order)
}

func (s *OrderService) validateProducts(items []*model.OrderItem) error {
	for _, item := range items {
		// Check if product exists
		product, err := s.productRepo.GetProductByID(item.ProductID)
		if err != nil || product == nil {
			return fmt.Errorf("product with ID %d does not exist", item.ProductID)
		}

		// TODO: implement inventory
		// Check inventory
		// if product.Inventory < item.Quantity {
		// 	return fmt.Errorf("not enough inventory for product ID %d", item.ProductID)
		// }

		// Check price
		if product.Price != item.Price {
			return fmt.Errorf("price mismatch for product ID %d: requested %.2f, actual %.2f", item.ProductID, item.Price, product.Price)
		}
	}
	return nil
}

func (s *OrderService) GetUserOrderItems(userID, orderID int) ([]*model.OrderItemWithProduct, error) {
	return s.orderRepo.GetUserOrderItems(userID, orderID)
}
