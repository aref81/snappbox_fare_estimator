package haversine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestHaversine tests the Haversine, the output should match the expected distance
func TestHaversine(t *testing.T) {
	// Test Case 1: Distance between New York (40.7128° N, 74.0060° W) and London (51.5074° N, 0.1278° W)
	lat1NY, lon1NY := 40.7128, -74.0060
	lat2LD, lon2LD := 51.5074, -0.1278
	expectedDistanceNYLD := 5570.0

	distanceNYLD := Haversine(lat1NY, lon1NY, lat2LD, lon2LD)
	assert.InDelta(t, expectedDistanceNYLD, distanceNYLD, 0.5, "Test Case 1: distance should be around 5570 km")

	// Test Case 2: Distance between Paris (48.8566° N, 2.3522° E) and Berlin (52.5200° N, 13.4050° E)
	lat1Paris, lon1Paris := 48.8566, 2.3522
	lat2Berlin, lon2Berlin := 52.5200, 13.4050
	expectedDistanceParisBerlin := 878.0

	distanceParisBerlin := Haversine(lat1Paris, lon1Paris, lat2Berlin, lon2Berlin)
	assert.InDelta(t, expectedDistanceParisBerlin, distanceParisBerlin, 1.0, "Test Case 2: distance should be around 878 km")

	// Test Case 3: Distance between the same point should be zero
	lat, lon := 35.7025, 51.4097
	expectedDistanceSamePoint := 0.0

	distanceSamePoint := Haversine(lat, lon, lat, lon)
	assert.Equal(t, expectedDistanceSamePoint, distanceSamePoint, "Distance between the same point should be 0")
}
