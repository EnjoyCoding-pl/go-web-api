package use_cases

import (
	"context"
	dto "go-web-api/features/user/app/models"
	"go-web-api/features/user/domain/models"
)

type GetByLoginStorage interface {
	GetByLogin(login string, ctx context.Context) (*models.User, error)
}

type TokenProvider interface {
	Generate(userId int) (*string, error)
}

type loginUserUseCase struct {
	storage       GetByLoginStorage
	tokenProvider TokenProvider
}

func NewLoginUserUseCase(s GetByLoginStorage, p TokenProvider) *loginUserUseCase {
	return &loginUserUseCase{storage: s, tokenProvider: p}
}

func (uc *loginUserUseCase) Execute(u *dto.LoginUserDto, ctx context.Context) (*string, error) {
	user, err := uc.storage.GetByLogin(u.Login, ctx)
	if err != nil {
		return nil, err
	}

	err = user.CheckPassword(u.Password)
	if err != nil {
		return nil, err
	}

	token, err := uc.tokenProvider.Generate(int(user.ID))

	if err != nil {
		return nil, err
	}

	return token, nil
}
