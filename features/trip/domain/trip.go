package domain

import (
	"time"
)

type Trip struct {
	Country string
	Begin   time.Time
	End     time.Time
}

func NewTrip() *Trip {
	return &Trip{}
}
