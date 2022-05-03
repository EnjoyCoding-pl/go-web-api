package use_cases

import (
	"context"
	dto "go-web-api/features/trip/app/models"
	"go-web-api/features/trip/domain/models"
)

type GetAllStorage interface {
	GetAll(sp *models.TripSearchParam, ctx context.Context) (*[]models.Trip, error)
}

type getAllUseCase struct {
	storage GetAllStorage
}

func NewGetAllUseCase(s GetAllStorage) *getAllUseCase {
	return &getAllUseCase{
		storage: s,
	}
}

func (u *getAllUseCase) Execute(sp *models.TripSearchParam, ctx context.Context) (*[]dto.TripDto, error) {

	result, err := u.storage.GetAll(sp, ctx)

	if err != nil {
		return nil, err
	}

	trips := make([]dto.TripDto, len(*result))

	for i, t := range *result {
		trips[i] = *dto.NewTripDto(t)
	}

	return &trips, nil
}
