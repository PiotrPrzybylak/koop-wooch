package memory

import (
	"github.com/PiotrPrzybylak/koop-wooch/domain"
	"strconv"
)

func NewSupplierRepository() domain.SupplierRepository {
	return &supplierRepository{}
}

type supplierRepository struct {
	nextID    int
	suppliers []domain.Supplier
}

func (r *supplierRepository) Create(supplier domain.Supplier) (string, error) {
	r.nextID++
	supplier.ID = strconv.Itoa(r.nextID)
	r.suppliers = append(r.suppliers, supplier)
	return supplier.ID, nil
}

func (s supplierRepository) ListAll() ([]domain.Supplier, error) {
	return s.suppliers, nil
}
