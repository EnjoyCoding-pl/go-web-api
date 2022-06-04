package use_cases

import (
	"context"
	dto "go-web-api/features/user/app/models"
	"go-web-api/features/user/domain/models"
)

type AddStorage interface {
	Add(u *models.User, ctx context.Context) error
}

type registerUserUseCase struct {
	storage AddStorage
}

func NewRegisterUserUseCase(s AddStorage) *registerUserUseCase {
	return &registerUserUseCase{storage: s}
}

func (uc *registerUserUseCase) Execute(u dto.RegisterUserDto, ctx context.Context) error {
	user, err := models.NewUser(u.Login, u.Password, u.RepeatedPassword)
	if err != nil {
		return err
	}
	return uc.storage.Add(user, ctx)
}
