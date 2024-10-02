package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestNewDelivery tests creating a new delivery with an ID
func TestNewDelivery(t *testing.T) {
	delivery := NewDelivery(1)
	assert.Equal(t, 1, delivery.ID, "Delivery ID should be 1")
	assert.Empty(t, delivery.Segments, "New delivery should have no segments")
}

// TestAddSegment tests adding two valid segments
func TestAddSegment(t *testing.T) {
	delivery := NewDelivery(1)

	// Define two valid DeliveryPoints (timestamps in seconds)
	startPoint := DeliveryPoint{
		DeliveryID: 1,
		Latitude:   35.700,
		Longitude:  51.400,
		Timestamp:  1000,
	}
	endPoint := DeliveryPoint{
		DeliveryID: 1,
		Latitude:   35.701,
		Longitude:  51.401,
		Timestamp:  1100,
	}

	// Add a valid segment
	err := delivery.AddSegment(startPoint, endPoint)
	assert.NoError(t, err, "Adding a valid segment should not produce an error")

	// Check that the segment was added
	assert.Len(t, delivery.Segments, 1, "Delivery should have one segment")

	segment := delivery.Segments[0]
	assert.Equal(t, startPoint.Timestamp, segment.StartTime, "Start time of the segment should match the start point timestamp")
	assert.Greater(t, segment.Distance, 0.0, "Distance should be calculated")
	assert.Greater(t, segment.Speed, 0.0, "Speed should be calculated")
}

// TestAddSegmentInvalidTimeDiff tests when the case that timestamp difference is zero
func TestAddSegmentInvalidTimeDiff(t *testing.T) {
	delivery := NewDelivery(1)
	startPoint := DeliveryPoint{
		DeliveryID: 1,
		Latitude:   35.700,
		Longitude:  51.400,
		Timestamp:  1000,
	}
	endPoint := DeliveryPoint{
		DeliveryID: 1,
		Latitude:   35.701,
		Longitude:  51.401,
		Timestamp:  1000, // Same timestamp
	}

	err := delivery.AddSegment(startPoint, endPoint)
	assert.Error(t, err, "Adding a segment with zero time difference should produce an error")
	assert.EqualError(t, err, "failed to calculate the speed, timeDiff = 0.000000", "Should return a time difference error")
}

// TestAddSegmentInvalidSpeed tests when speed exceeds the limit (hence the segment is invalid)
func TestAddSegmentInvalidSpeed(t *testing.T) {
	delivery := NewDelivery(1)
	startPoint := DeliveryPoint{
		DeliveryID: 1,
		Latitude:   35.700,
		Longitude:  51.400,
		Timestamp:  1000,
	}
	endPoint := DeliveryPoint{
		DeliveryID: 1,
		Latitude:   36.700, // Speed ~ 5143 >> 100
		Longitude:  52.400,
		Timestamp:  1100,
	}

	err := delivery.AddSegment(startPoint, endPoint)
	assert.Error(t, err, "Adding a segment with unrealistic speed should produce an error")
	assert.Contains(t, err.Error(), "invalid processor point, speed =", "Should return a speed error")
}

// TestValidateSegment tests the segment validation mechanism
func TestValidateSegment(t *testing.T) {
	// valid segment
	segment := DeliverySegment{
		Speed: 90.0,
	}
	err := validateSegment(segment)
	assert.NoError(t, err, "Valid segment speed should not return an error")

	// invalid segment (speed > 100)
	segment.Speed = 120.0
	err = validateSegment(segment)
	assert.Error(t, err, "Invalid segment speed should return an error")
	assert.Contains(t, err.Error(), "invalid processor point, speed =", "Error message should contain speed validation issue")
}

func TestCalculateSpeedAndDistance(t *testing.T) {
	startPoint := DeliveryPoint{
		DeliveryID: 1,
		Latitude:   35.700,
		Longitude:  51.400,
		Timestamp:  1000,
	}
	endPoint := DeliveryPoint{
		DeliveryID: 1,
		Latitude:   35.701,
		Longitude:  51.401,
		Timestamp:  1100,
	}

	speed, distance, err := calculateSpeedAndDistance(startPoint, endPoint)
	expectedSpeed, expectedDistance := 5.15, 0.14
	assert.NoError(t, err, "Calculating speed and distance should not produce an error")
	assert.Greater(t, distance, 0.0, "Distance should be greater than 0")
	assert.Greater(t, speed, 0.0, "Speed should be greater than 0")
	assert.InDelta(t, expectedSpeed, speed, 0.5, "speed should be around 5.15")
	assert.InDelta(t, expectedDistance, distance, 0.5, "distance should be around 0.14")

	// Test with invalid time difference
	startPoint.Timestamp = 1000
	endPoint.Timestamp = 1000 // Same timestamp

	speed, distance, err = calculateSpeedAndDistance(startPoint, endPoint)
	assert.Error(t, err, "Should return an error for zero time difference")
	assert.EqualError(t, err, "failed to calculate the speed, timeDiff = 0.000000", "Should return a time difference error")
}
