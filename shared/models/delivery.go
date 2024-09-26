package models

import (
	"fmt"
	"github.com/aref81/snappbox_fare_estimator/shared/haversine"
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
func (d *Delivery) AddPoint(lat, lng float64, timestamp int64) error {
	newPoint := DeliveryPoint{
		Latitude:  lat,
		Longitude: lng,
		Timestamp: timestamp,
	}

	if len(d.Points) != 0 {
		err := validateDeliveryPoints(d.Points[len(d.Points)-1], newPoint)
		if err != nil {
			return err
		}
	}

	d.Points = append(d.Points, newPoint)
	return nil
}

// NewDelivery initializes a new Delivery
func NewDelivery(id int) *Delivery {
	return &Delivery{
		ID:     id,
		Points: []DeliveryPoint{},
	}
}

// ValidateDeliveryPoints checks the validity of two consecutive delivery points
func validateDeliveryPoints(p1 DeliveryPoint, p2 DeliveryPoint) error {
	timeDiff := float64(p2.Timestamp-p1.Timestamp) / 3600.0
	if timeDiff == 0 {
		// skipping zero time differences
		return nil
	}

	distance := haversine.Haversine(p1.Latitude, p1.Longitude, p2.Latitude, p2.Longitude)
	speed := distance / timeDiff

	if speed <= 100.0 {
		return nil
	} else {
		return fmt.Errorf("invalid delivery point, speed = %f", speed)
	}
}
