package models

import (
	"testing"
	"time"
)

func TestGetBeginDate_Success(t *testing.T) {
	now := time.Now()
	trip := &Trip{
		Points: []TripPoint{
			{
				Begin: now,
				End:   now.AddDate(0, 0, 1),
			},
			{
				Begin: now.AddDate(0, 0, 1),
				End:   now.AddDate(0, 0, 2),
			},
		},
	}

	bd := trip.GetBeginDate()

	if !now.Equal(bd) {
		t.Fatalf("Begin date isn't equal to now, Begin date: %v, Now: %v", bd, now)
	}
}

func TestGetEndDate_Success(t *testing.T) {
	now := time.Now()
	trip := &Trip{
		Points: []TripPoint{
			{
				Begin: now,
				End:   now.Add(1),
			},
			{
				Begin: now.Add(1),
				End:   now.Add(2),
			},
		},
	}
	bd := trip.GetEndDate()
	expected := now.Add(2)

	if !expected.Equal(bd) {
		t.Fatalf("End date isn't equal to now + two days, End date: %v, Expected: %v", bd, expected)
	}
}
