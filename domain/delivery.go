package domain

type Delivery struct {
	Supplier string
	Category string
	Price    float64
	Unit     string
	Quantity float64
}

type DeliveryService interface {
	Create(delivery Delivery) (string, error)
	All() ([]Delivery, error)
}

type DeliveryRepository interface {
	Create(delivery Delivery) (string, error)
	All() ([]Delivery, error)
}

func NewDeliverytService(repo DeliveryRepository) DeliveryService {
	return deliveryService{repo}
}

type deliveryService struct {
	repo DeliveryRepository
}

func (s deliveryService) Create(delivery Delivery) (string, error) {
	return s.repo.Create(delivery)
}

func (s deliveryService) All() ([]Delivery, error) {
	return s.repo.All()
}
