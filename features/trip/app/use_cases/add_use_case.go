package use_cases

import (
	"context"
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

func (u *addUseCase) Execute(t models.Trip, ctx context.Context) error {
	return u.storage.Add(t, ctx)
}
