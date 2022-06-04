package storages

import (
	"context"
	"errors"
	"go-web-api/features/user/domain/models"
	"go-web-api/internal/globals"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type userPostgresStorage struct {
	db *gorm.DB
}

func NewUserPostgresStorage(db *gorm.DB) *userPostgresStorage {
	return &userPostgresStorage{db: db}
}

func (s *userPostgresStorage) Add(u *models.User, ctx context.Context) error {
	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(ctx, "postgres-add-query")
	defer span.End()

	tx := s.db.WithContext(spanCtx).Create(u)
	return tx.Error
}

func (s *userPostgresStorage) GetByLogin(login string, ctx context.Context) (*models.User, error) {
	spanCtx, span := otel.Tracer(globals.TracerAppName).Start(ctx, "postgres-get-by-login-query")
	defer span.End()

	var user models.User

	tx := s.db.WithContext(spanCtx).Where("login = ?", login).First(&user)

	if tx.Error != nil {
		return nil, tx.Error
	}

	if tx.RowsAffected <= 0 {
		return nil, errors.New("missing user")
	}

	return &user, nil
}
