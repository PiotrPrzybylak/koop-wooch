package domain

type Product struct {
	ID       string
	Name     string
	Category string
	Price    float64
	Unit     string
	Quantity float64
}

type ProductService interface {
	Create(product Product) (string, error)
	All() ([]Product, error)
}

type ProductRepository interface {
	Create(product Product) (string, error)
	All() ([]Product, error)
}

func NewProductService(repo ProductRepository) ProductService {
	return productService{repo}
}

type productService struct {
	repo ProductRepository
}

func (s productService) Create(product Product) (string, error) {
	return s.repo.Create(product)
}

func (s productService) All() ([]Product, error) {
	return s.repo.All()
}
