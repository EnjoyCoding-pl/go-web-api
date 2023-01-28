package storages

import (
	"context"
	"go-web-api/features/trip/domain/models"
	"testing"

	"go-web-api/internal/containers"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestPostgresStorage(t *testing.T) {

	ctx := context.Background()

	database := "trip_db"

	container, err := containers.NewPostgresContainer(database, ctx)
	if err != nil {
		panic(err)
	}

	db := container.Db

	err = db.AutoMigrate(models.NewTrip(), models.NewTripPoint())

	if err != nil {
		panic(err)
	}

	storage := NewPostgresStorage(db)

	t.Cleanup(func() {
		container.Terminate(ctx)
	})

	t.Run("Add success", func(t *testing.T) {
		newTrip := models.NewTrip()
		newTrip.Country = "Add test"

		err := storage.Add(*newTrip, ctx)
		assert.NoError(t, err)
	})

	t.Run("Get all success", func(t *testing.T) {

		newTrip := models.NewTrip()
		newTrip.Country = "Get all"

		tx := db.Create(&newTrip)
		assert.NoError(t, tx.Error)

		trips, err := storage.GetAll(models.NewTripSearchParam(1, 10), ctx)

		assert.NoError(t, err)
		assert.Greater(t, len(*trips), 0)
	})

	t.Run("Get by id success", func(t *testing.T) {

		newTrip := models.NewTrip()
		newTrip.Country = "Get by id"

		tx := db.Create(&newTrip)
		assert.NoError(t, tx.Error)

		trip, err := storage.Get(int(newTrip.ID), ctx)

		assert.NoError(t, err)
		assert.NotEmpty(t, trip)
	})

	t.Run("Delete success", func(t *testing.T) {

		newTrip := models.NewTrip()
		newTrip.Country = "Get by id"

		tx := db.Create(&newTrip)
		assert.NoError(t, tx.Error)

		err := storage.Delete(int(newTrip.ID), ctx)

		assert.NoError(t, err)

		tx = db.First(models.NewTrip(), newTrip.ID)
		assert.ErrorIs(t, tx.Error, gorm.ErrRecordNotFound)
	})
}
