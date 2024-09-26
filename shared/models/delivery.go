package models

import (
	"fmt"
)

// DeliveryPoint represents a single GPS coordination for a Delivery, each Delivery contains a sorted list of
// these points
type DeliveryPoint struct {
	Latitude  float64
	Longitude float64
	Timestamp int64
}

// Delivery represents an individual Delivery Data, which includes and ID and multiple Delivery Points
type Delivery struct {
	ID     int
	Points []DeliveryPoint
}

// AddPoint adds a new DeliveryPoint to the Delivery
func (d *Delivery) AddPoint(lat, lng float64, timestamp int64) {
	d.Points = append(d.Points, DeliveryPoint{
		Latitude:  lat,
		Longitude: lng,
		Timestamp: timestamp,
	})
}

// NewDelivery initializes a new Delivery
func NewDelivery(id int) *Delivery {
	return &Delivery{
		ID:     id,
		Points: []DeliveryPoint{},
	}
}

// ValidateDelivery checks the validity of a delivery (just a basic validation as an example of practice)
func ValidateDelivery(delivery *Delivery) error {
	if len(delivery.Points) == 0 {
		return fmt.Errorf("delivery %d has no points", delivery.ID)
	}

	// Check if the points have valid latitude and longitude
	for _, point := range delivery.Points {
		if point.Latitude < -90 || point.Latitude > 90 {
			return fmt.Errorf("invalid latitude in delivery %d", delivery.ID)
		}
		if point.Longitude < -180 || point.Longitude > 180 {
			return fmt.Errorf("invalid longitude in delivery %d", delivery.ID)
		}
	}

	return nil
}
