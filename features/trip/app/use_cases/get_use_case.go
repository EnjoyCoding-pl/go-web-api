package use_cases

import (
	"context"
	dto "go-web-api/features/trip/app/models"
	"go-web-api/features/trip/domain/models"
)

type GetStorage interface {
	Get(id int, ctx context.Context) (models.Trip, error)
}

type getUseCase struct {
	storage GetStorage
}

func NewGetUseCase(s GetStorage) *getUseCase {
	return &getUseCase{storage: s}
}

func (u *getUseCase) Execute(id int, ctx context.Context) (*dto.TripDto, error) {
	t, err := u.storage.Get(id, ctx)

	if err != nil {
		return nil, err
	}

	return dto.NewTripDto(t), nil
}
