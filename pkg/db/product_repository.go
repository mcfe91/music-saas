package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"music-saas/pkg/model"
)

type PostgresProductRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProductRepository(db *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (repo *PostgresProductRepository) CreateProduct(product *model.Product) error {
	_, err := repo.db.Exec(context.Background(),
		"INSERT INTO products (name, description, price) VALUES ($1, $2, $3)",
		product.Name, product.Description, product.Price)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresProductRepository) GetProductByID(id int) (*model.Product, error) {
	product := &model.Product{}
	row := repo.db.QueryRow(context.Background(),
		"SELECT id, name, description, price, created_at FROM products WHERE id = $1", id)
	err := row.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return product, nil
}

func (repo *PostgresProductRepository) UpdateProduct(product *model.Product) error {
	_, err := repo.db.Exec(context.Background(),
		"UPDATE products SET name = $1, description = $2, price = $3 WHERE id = $4",
		product.Name, product.Description, product.Price, product.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresProductRepository) DeleteProduct(id int) error {
	_, err := repo.db.Exec(context.Background(),
		"DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PostgresProductRepository) GetProducts(limit, offset int) ([]*model.Product, error) {
	rows, err := repo.db.Query(context.Background(),
		"SELECT id, name, description, price, created_at FROM products ORDER BY id LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.CreatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}
