package service

import (
	"music-saas/pkg/db"
	"music-saas/pkg/model"
)

type ProductService struct {
	productRepo *db.PostgresProductRepository
}

func NewProductService(repo *db.PostgresProductRepository) *ProductService {
	return &ProductService{productRepo: repo}
}

func (s *ProductService) CreateProduct(product *model.Product) error {
	return s.productRepo.CreateProduct(product)
}

func (s *ProductService) GetProductByID(id int) (*model.Product, error) {
	return s.productRepo.GetProductByID(id)
}

func (s *ProductService) UpdateProduct(product *model.Product) error {
	return s.productRepo.UpdateProduct(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.productRepo.DeleteProduct(id)
}

func (s *ProductService) GetProducts(limit, offset int) ([]*model.Product, error) {
	return s.productRepo.GetProducts(limit, offset)
}
