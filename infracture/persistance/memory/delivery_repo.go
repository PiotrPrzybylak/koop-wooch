package memory

import (
	"github.com/PiotrPrzybylak/koop-wooch/domain"
)

func NewDeliveryRepository() domain.DeliveryRepository {
	return &deliveryRepository{}
}

type deliveryRepository struct {
	deliverys []domain.Delivery
}

func (r *deliveryRepository) Create(delivery domain.Delivery) (string, error) {
	r.deliverys = append(r.deliverys, delivery)
	return delivery.Category, nil
}

func (s deliveryRepository) All() ([]domain.Delivery, error) {
	return s.deliverys, nil
}
