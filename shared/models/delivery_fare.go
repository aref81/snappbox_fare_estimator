package models

// DeliveryFare holds Fare amount calculated for delivery
type DeliveryFare struct {
	ID   int
	Fare float64
}

// NewDeliveryFare initializes a new DeliveryFare
func NewDeliveryFare(id int, fare float64) *DeliveryFare {
	return &DeliveryFare{
		ID:   id,
		Fare: fare,
	}
}
