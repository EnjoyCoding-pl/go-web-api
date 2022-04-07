package models

import (
	"time"

	"gorm.io/gorm"
)

type Trip struct {
	gorm.Model
	Country string
	Points  []TripPoint
	Begin   time.Time
	End     time.Time
}

func NewTrip() *Trip {
	return &Trip{}
}
