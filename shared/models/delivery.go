package models

import (
	"fmt"
	"github.com/aref81/snappbox_fare_estimator/shared/haversine"
)

// DeliveryPoint represents a single GPS coordination for a Delivery
type DeliveryPoint struct {
	Latitude  float64
	Longitude float64
	Timestamp int64
}

// DeliverySegment represents a segment of the road traveled, including two DeliveryPoints and Speed calculated for it.
type DeliverySegment struct {
	StartPoint DeliveryPoint
	EndPoint   DeliveryPoint
	Speed      float64
}

// Delivery represents an individual Delivery Data, which includes and ID and multiple DeliverySegments
type Delivery struct {
	ID       int
	Segments []DeliverySegment
}

// AddSegment adds a new DeliverySegment to the Delivery after validation
func (d *Delivery) AddSegment(startPoint DeliveryPoint, endPoint DeliveryPoint) error {
	segmentSpeed, err := calculateSpeed(startPoint, endPoint)
	if err != nil {
		return err
	}

	segment := DeliverySegment{
		StartPoint: startPoint,
		EndPoint:   endPoint,
		Speed:      segmentSpeed,
	}

	err = validateSegment(segment)
	if err != nil {
		return err
	}

	d.Segments = append(d.Segments, segment)
	return nil
}

// NewDelivery initializes a new Delivery
func NewDelivery(id int) *Delivery {
	return &Delivery{
		ID:       id,
		Segments: []DeliverySegment{},
	}
}

// ValidateDeliverySegment checks the validity of a segment by calculating the speed of it
func validateSegment(segment DeliverySegment) error {
	if segment.Speed <= 100.0 {
		return nil
	} else {
		return fmt.Errorf("invalid delivery point, speed = %f", segment.Speed)
	}
}

// calculateSpeed calculates the speed for a segment using haversine distance
func calculateSpeed(p1 DeliveryPoint, p2 DeliveryPoint) (float64, error) {
	timeDiff := float64(p2.Timestamp-p1.Timestamp) / 3600.0
	if timeDiff == 0 {
		// skipping zero time differences
		return 0, fmt.Errorf("failed to calculate the speed, timeDiff = %f", timeDiff)
	}

	distance := haversine.Haversine(p1.Latitude, p1.Longitude, p2.Latitude, p2.Longitude)
	speed := distance / timeDiff
	return speed, nil
}
