package use_cases

import (
	"context"
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

func (u *getUseCase) Execute(id int, ctx context.Context) (models.Trip, error) {
	return u.storage.Get(id, ctx)
}
