package use_cases

import (
	"context"
	dto "go-web-api/features/trip/app/models"
	"go-web-api/features/trip/domain/models"
)

type UpdateStorage interface {
	Update(t models.Trip, ctx context.Context) error
}

type updateUseCase struct {
	storage UpdateStorage
}

func NewUpdatetUseCase(s UpdateStorage) *updateUseCase {
	return &updateUseCase{storage: s}
}

func (u *updateUseCase) Execute(t dto.TripDto, ctx context.Context) error {
	return u.storage.Update(*t.ToDomain(), ctx)
}
