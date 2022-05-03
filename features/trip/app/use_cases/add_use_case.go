package use_cases

import (
	"context"
	dto "go-web-api/features/trip/app/models"
	"go-web-api/features/trip/domain/models"
)

type AddStorage interface {
	Add(t models.Trip, ctx context.Context) error
}

type addUseCase struct {
	storage AddStorage
}

func NewAddUseCase(s AddStorage) *addUseCase {
	return &addUseCase{storage: s}
}

func (u *addUseCase) Execute(t dto.TripDto, ctx context.Context) error {
	return u.storage.Add(*t.ToDomain(), ctx)
}
