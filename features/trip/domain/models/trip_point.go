package models

import "gorm.io/gorm"

type TripPoint struct {
	gorm.Model
	Place  string
	TripID uint
}

func NewTripPoint() *TripPoint {
	return &TripPoint{}
}
