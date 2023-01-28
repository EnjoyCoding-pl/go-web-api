package storages

import (
	"context"
	"go-web-api/features/user/domain/models"
	"go-web-api/internal/containers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresStorage(t *testing.T) {

	ctx := context.Background()

	database := "user_db"

	container, err := containers.NewPostgresContainer(database, ctx)
	assert.NoError(t, err)

	db := container.Db

	err = db.AutoMigrate(&models.User{})

	assert.NoError(t, err)

	storage := NewUserPostgresStorage(db)

	t.Cleanup(func() {
		container.Terminate(ctx)
	})

	t.Run("Add success", func(t *testing.T) {
		u, err := models.NewUser("login", "password", "password")
		assert.NoError(t, err)
		storage.Add(u, ctx)
	})

	t.Run("Get by login success", func(t *testing.T) {
		u, err := models.NewUser("login-2", "password", "password")
		assert.NoError(t, err)
		tx := db.Create(u)
		assert.NoError(t, tx.Error)

		user, err := storage.GetByLogin("login-2", ctx)

		assert.NoError(t, err)
		assert.NotEmpty(t, user)
	})
}
