package storages

import (
	"context"
	"go-web-api/features/trip/domain/models"

	"gorm.io/gorm"
)

type postgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage(db *gorm.DB) *postgresStorage {
	return &postgresStorage{db: db}
}

func (s *postgresStorage) Add(t models.Trip, ctx context.Context) error {
	tx := s.db.Create(&t)

	return tx.Error
}

func (s *postgresStorage) Update(t models.Trip, ctx context.Context) error {

	var points []models.TripPoint
	s.db.Where("trip_id = ?", t.ID).Find(&points)

	deleted := make([]uint, 0)

	for _, point := range points {
		exists := false
		for _, newPoint := range t.Points {
			if newPoint.ID == point.ID {
				exists = true
				break
			}
		}
		if !exists {
			deleted = append(deleted, point.ID)
		}
	}

	deleteTx := s.db.Delete(models.NewTripPoint(), deleted)

	if deleteTx.Error != nil {
		return deleteTx.Error
	}

	saveTx := s.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(&t)

	return saveTx.Error
}

func (s *postgresStorage) Delete(id int, ctx context.Context) error {
	tx := s.db.Delete(models.NewTrip(), id)

	return tx.Error
}

func (s *postgresStorage) Get(id int, ctx context.Context) (models.Trip, error) {
	t := models.NewTrip()

	tx := s.db.Preload("Points").First(t, id)

	return *t, tx.Error
}
