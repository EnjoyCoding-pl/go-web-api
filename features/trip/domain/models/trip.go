package models

import (
	"math"
	"time"

	"gorm.io/gorm"
)

type Trip struct {
	gorm.Model
	Country string
	Points  []TripPoint
}

func NewTrip() *Trip {
	return &Trip{}
}

func (t *Trip) GetBeginDate() time.Time {

	var min time.Time
	if t.Points != nil {
		for _, tp := range t.Points {
			if min.IsZero() || tp.Begin.Before(min) {
				min = tp.Begin
			}
		}
	}

	return min
}

func (t *Trip) GetEndDate() time.Time {
	var max time.Time
	if t.Points != nil {
		for _, tp := range t.Points {
			if max.IsZero() || tp.End.After(max) {
				max = tp.End
			}
		}
	}
	return max

}

func (t *Trip) GetTotalDistance() float64 {
	var sum float64 = 0
	if t.Points != nil {
		for i, tp := range t.Points {
			if i > 0 {
				sum += calculateDistance(t.Points[i-1], tp)
			}
		}
	}
	return sum
}

func calculateDistance(leftPoint TripPoint, rightPoint TripPoint) float64 {
	earthRediusKm := 6371

	deltaLat := degreesToRadians(rightPoint.Latitude - leftPoint.Latitude)
	deltaLon := degreesToRadians(rightPoint.Longitude - leftPoint.Longitude)

	leftLat := degreesToRadians(leftPoint.Latitude)
	rightLat := degreesToRadians(rightPoint.Latitude)

	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Pow(math.Sin(deltaLon/2), 2)*math.Cos(leftLat)*math.Cos(rightLat)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return float64(earthRediusKm) * c
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
