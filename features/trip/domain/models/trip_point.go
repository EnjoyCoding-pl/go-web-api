package models

import (
	"time"

	"gorm.io/gorm"
)

type TripPoint struct {
	gorm.Model
	Place     string
	Latitude  float64
	Longitude float64
	Begin     time.Time
	End       time.Time
	TripID    uint
}

func NewTripPoint() *TripPoint {
	return &TripPoint{}
}
