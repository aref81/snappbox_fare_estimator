package processor

import (
	"testing"
	"time"

	"github.com/aref81/snappbox_fare_estimator/atalanta/config"
	"github.com/aref81/snappbox_fare_estimator/shared/models"
	"github.com/stretchr/testify/assert"
)

func TestCalculateFare_MovingDayTime(t *testing.T) {
	// Define the fare rules
	fareRules := config.FareRulesConfig{
		FlagAmount:           5.0,
		MovingDayFarePerKm:   10.0,
		MovingNightFarePerKm: 15.0,
		IdleFarePerHour:      2.0,
		MinFare:              20.0,
	}
	timeBoundaries := config.TimeBoundariesConfig{
		DayStartHour: 6,
		DayEndHour:   20,
	}

	calculator := &fareCalculator{
		fareConfig:     fareRules,
		timeBoundaries: timeBoundaries,
	}

	delivery := &models.Delivery{
		ID: 1,
		Segments: []models.DeliverySegment{
			{
				StartTime:   time.Date(2023, 9, 30, 10, 0, 0, 0, time.UTC).Unix(), // 10:00 AM (daytime)
				ElapsedTime: 0.5,
				Distance:    5.0,
				Speed:       50.0,
			},
		},
	}

	fare := calculator.calculateFare(delivery)
	expectedFare := fareRules.FlagAmount + (5.0 * fareRules.MovingDayFarePerKm)
	assert.Equal(t, expectedFare, fare, "The fare for a moving segment during the day should be correctly calculated")
}

func TestCalculateFare_MovingNightTime(t *testing.T) {
	fareRules := config.FareRulesConfig{
		FlagAmount:           5.0,
		MovingDayFarePerKm:   10.0,
		MovingNightFarePerKm: 15.0,
		IdleFarePerHour:      2.0,
		MinFare:              20.0,
	}
	timeBoundaries := config.TimeBoundariesConfig{
		DayStartHour: 6,
		DayEndHour:   20,
	}

	calculator := &fareCalculator{
		fareConfig:     fareRules,
		timeBoundaries: timeBoundaries,
	}

	// Define a delivery with a moving segment during the night
	delivery := &models.Delivery{
		ID: 1,
		Segments: []models.DeliverySegment{
			{
				StartTime:   time.Date(2023, 9, 30, 22, 0, 0, 0, time.UTC).Unix(), // 10:00 PM (nighttime)
				ElapsedTime: 0.5,
				Distance:    5.0,
				Speed:       50.0,
			},
		},
	}

	fare := calculator.calculateFare(delivery)
	expectedFare := fareRules.FlagAmount + (5.0 * fareRules.MovingNightFarePerKm)
	assert.Equal(t, expectedFare, fare, "The fare for a moving segment during the night should be correctly calculated")
}

func TestCalculateFare_IdleSegment(t *testing.T) {
	fareRules := config.FareRulesConfig{
		FlagAmount:           5.0,
		MovingDayFarePerKm:   10.0,
		MovingNightFarePerKm: 15.0,
		IdleFarePerHour:      2.0,
		MinFare:              0.0,
	}
	timeBoundaries := config.TimeBoundariesConfig{
		DayStartHour: 6,
		DayEndHour:   20,
	}

	calculator := &fareCalculator{
		fareConfig:     fareRules,
		timeBoundaries: timeBoundaries,
	}

	// Define a delivery with an idle segment
	delivery := &models.Delivery{
		ID: 1,
		Segments: []models.DeliverySegment{
			{
				StartTime:   time.Date(2023, 9, 30, 10, 0, 0, 0, time.UTC).Unix(), // 10:00 AM (daytime)
				ElapsedTime: 1.0,
				Distance:    0.0,
				Speed:       0.0,
			},
		},
	}

	fare := calculator.calculateFare(delivery)
	expectedFare := fareRules.FlagAmount + (1.0 * fareRules.IdleFarePerHour)
	assert.Equal(t, expectedFare, fare, "The fare for an idle segment should be correctly calculated")
}

func TestCalculateFare_MinimumFare(t *testing.T) {
	fareRules := config.FareRulesConfig{
		FlagAmount:           5.0,
		MovingDayFarePerKm:   10.0,
		MovingNightFarePerKm: 15.0,
		IdleFarePerHour:      2.0,
		MinFare:              20.0,
	}
	timeBoundaries := config.TimeBoundariesConfig{
		DayStartHour: 6,
		DayEndHour:   20,
	}

	calculator := &fareCalculator{
		fareConfig:     fareRules,
		timeBoundaries: timeBoundaries,
	}

	// Define a delivery where the calculated fare is below the minimum fare
	delivery := &models.Delivery{
		ID: 1,
		Segments: []models.DeliverySegment{
			{
				StartTime:   time.Date(2023, 9, 30, 10, 0, 0, 0, time.UTC).Unix(), // 10:00 AM (daytime)
				ElapsedTime: 0.5,
				Distance:    0.5,
				Speed:       10.0,
			},
		},
	}

	fare := calculator.calculateFare(delivery)
	assert.Equal(t, fareRules.MinFare, fare, "The fare should be set to the minimum fare")
}

func TestCalculateFare_ConsecutiveSegments_DayAndNight(t *testing.T) {
	fareRules := config.FareRulesConfig{
		FlagAmount:           5.0,
		MovingDayFarePerKm:   10.0,
		MovingNightFarePerKm: 15.0,
		IdleFarePerHour:      2.0,
		MinFare:              20.0,
	}
	timeBoundaries := config.TimeBoundariesConfig{
		DayStartHour: 6,
		DayEndHour:   20,
	}

	calculator := &fareCalculator{
		fareConfig:     fareRules,
		timeBoundaries: timeBoundaries,
	}

	// - First segment: Starts during the day and ends right before nighttime
	// - Second segment: Starts immediately at nighttime
	segment1Start := time.Date(2023, 9, 30, 19, 30, 0, 0, time.UTC).Unix() // 7:30 PM (daytime)
	segment2Start := time.Date(2023, 9, 30, 20, 0, 0, 0, time.UTC).Unix()  // 8:00 PM (nighttime)

	delivery := &models.Delivery{
		ID: 1,
		Segments: []models.DeliverySegment{
			{
				StartTime:   segment1Start,
				ElapsedTime: 0.5,
				Distance:    5.0,
				Speed:       50.0,
			},
			{
				StartTime:   segment2Start,
				ElapsedTime: 0.5,
				Distance:    3.0,
				Speed:       50.0,
			},
		},
	}

	fare := calculator.calculateFare(delivery)
	expectedFare := fareRules.FlagAmount + (5.0 * fareRules.MovingDayFarePerKm) + (3.0 * fareRules.MovingNightFarePerKm)
	assert.Equal(t, expectedFare, fare, "The fare for a delivery with consecutive day and night segments should be correctly calculated")
}

func TestIsDayTime(t *testing.T) {
	timeBoundaries := config.TimeBoundariesConfig{
		DayStartHour: 6,
		DayEndHour:   20,
	}

	calculator := &fareCalculator{
		timeBoundaries: timeBoundaries,
	}

	dayTimestamp := time.Date(2023, 9, 30, 12, 0, 0, 0, time.UTC).Unix() // 12:00 PM (daytime)
	assert.True(t, calculator.isDayTime(dayTimestamp), "12:00 PM should be considered daytime")

	nightTimestamp := time.Date(2023, 9, 30, 22, 0, 0, 0, time.UTC).Unix() // 10:00 PM (nighttime)
	assert.False(t, calculator.isDayTime(nightTimestamp), "10:00 PM should be considered nighttime")
}
