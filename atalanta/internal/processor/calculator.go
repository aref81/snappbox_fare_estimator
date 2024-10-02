package processor

import (
	"github.com/aref81/snappbox_fare_estimator/atalanta/config"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"time"
)

type fareCalculator struct {
	fareConfig     config.FareRulesConfig
	timeBoundaries config.TimeBoundariesConfig
}

// CalculateFare calculates the fare amount for each processor based on fare rules
func (c *fareCalculator) calculateFare(delivery *models.Delivery) float64 {
	totalFare := c.fareConfig.FlagAmount

	for _, segment := range delivery.Segments {
		// Decide if the status is moving or idle
		if segment.Speed > 10 {
			// Determine if it's day or night fare
			if c.isDayTime(segment.StartTime) {
				totalFare += segment.Distance * c.fareConfig.MovingDayFarePerKm
			} else {
				totalFare += segment.Distance * c.fareConfig.MovingNightFarePerKm
			}
		} else {
			totalFare += c.fareConfig.IdleFarePerHour * segment.ElapsedTime
		}
	}

	// Check for minimum fare
	if totalFare < c.fareConfig.MinFare {
		totalFare = c.fareConfig.MinFare
	}

	return totalFare
}

func (c *fareCalculator) isDayTime(timestamp int64) bool {
	t := time.Unix(timestamp, 0).UTC()
	hour := t.Hour()
	return hour >= c.timeBoundaries.DayStartHour && hour < c.timeBoundaries.DayEndHour
}
