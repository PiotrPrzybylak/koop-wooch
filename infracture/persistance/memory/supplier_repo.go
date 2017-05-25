package memory

import (
	"errors"
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

func (r supplierRepository) ListAll() ([]domain.Supplier, error) {
	return r.suppliers, nil
}

func (r *supplierRepository) Delete(id string) error {
	var result int = -1
	for index, value := range r.suppliers {
		if value.ID == id {
			result = index
			break
		}

	}
	if result != -1 {
		r.suppliers = append(r.suppliers[:result], r.suppliers[result+1:]...)
		return nil
	} else {
		return errors.New("There is no supplier with given id " + id)
	}

}
