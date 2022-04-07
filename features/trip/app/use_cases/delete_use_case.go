package use_cases

import (
	"context"
)

type DeleteStorage interface {
	Delete(id int, ctx context.Context) error
}

type deleteUseCase struct {
	storage DeleteStorage
}

func NewDeleteUseCase(s DeleteStorage) *deleteUseCase {
	return &deleteUseCase{storage: s}
}

func (u *deleteUseCase) Execute(id int, ctx context.Context) error {
	return u.storage.Delete(id, ctx)
}
