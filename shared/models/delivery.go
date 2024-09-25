package models

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
