package memory

import (
	"strconv"
	"github.com/PiotrPrzybylak/koop-wooch/domain"
)

func NewProductRepository() domain.ProductRepository {
	return &productRepository{}
}

type productRepository struct {
	nextID   int
	products []domain.Product
}

func (r *productRepository) Create(product domain.Product) (string, error) {
	r.nextID++
	product.ID = strconv.Itoa(r.nextID)
	r.products = append(r.products, product)
	return product.ID, nil
}

func (s productRepository) All() ([]domain.Product, error) {
	return s.products, nil
}
