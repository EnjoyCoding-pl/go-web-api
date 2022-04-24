package use_cases

import (
	"context"
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

func (u *getAllUseCase) Execute(sp *models.TripSearchParam, ctx context.Context) (*[]models.Trip, error) {
	return u.storage.GetAll(sp, ctx)

}
