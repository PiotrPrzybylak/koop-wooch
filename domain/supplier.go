package domain

import (
	"time"
)

type Supplier struct {
	ID          string
	Name        string
	DeliveryDay time.Weekday
}

type SupplierService interface {
	Create(supplier Supplier) (string, error)
	ListAll() ([]Supplier, error)
	Delete(id string) (error)
}

type SupplierRepository interface {
	Create(supplier Supplier) (string, error)
	ListAll() ([]Supplier, error)
	Delete(id string) (error)
}

func NewSupplierService(repo SupplierRepository) SupplierService {
	return supplierService{repo}
}

type supplierService struct {
	repo SupplierRepository
}

func (s supplierService) Create(supplier Supplier) (string, error) {
	return s.repo.Create(supplier)
}

func (s supplierService) ListAll() ([]Supplier, error) {
	return s.repo.ListAll()
}

func (s supplierService) Delete(id string) (error) {
	return s.repo.Delete(id)
}
