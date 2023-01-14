package models

import (
	"go-web-api/features/trip/domain/models"
	"time"

	"gorm.io/gorm"
)

type TripDto struct {
	Id            uint
	Country       string
	Points        []TripPointDto
	Begin         time.Time
	End           time.Time
	TotalDistance float64
}
type TripPointDto struct {
	Id        uint
	Place     string
	Latitude  float64
	Longitude float64
	Begin     time.Time
	End       time.Time
}

func NewTripDto(t models.Trip) *TripDto {
	dto := &TripDto{
		Id:            t.ID,
		Country:       t.Country,
		Begin:         t.GetBeginDate(),
		End:           t.GetEndDate(),
		TotalDistance: t.GetTotalDistance(),
	}
	if t.Points != nil {

		for _, tp := range t.Points {
			dto.Points = append(dto.Points, TripPointDto{
				Id:        tp.ID,
				Place:     tp.Place,
				Latitude:  tp.Latitude,
				Longitude: tp.Longitude,
				Begin:     tp.Begin,
				End:       tp.End,
			})
		}
	}

	return dto
}

func (t *TripDto) ToDomain() *models.Trip {
	domain := &models.Trip{
		Model: gorm.Model{
			ID: t.Id,
		},
		Country: t.Country,
	}
	if t.Points != nil {
		for _, tp := range t.Points {
			domain.Points = append(domain.Points, models.TripPoint{
				Model: gorm.Model{
					ID: tp.Id,
				},
				TripID:    t.Id,
				Place:     tp.Place,
				Latitude:  tp.Latitude,
				Longitude: tp.Longitude,
				Begin:     tp.Begin,
				End:       tp.End,
			})
		}
	}

	return domain
}
