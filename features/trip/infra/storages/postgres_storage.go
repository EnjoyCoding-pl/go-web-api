package storages

import (
	"context"
	"go-web-api/features/trip/domain/models"
	"go-web-api/internal/globals"

	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
)

type postgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage(db *gorm.DB) *postgresStorage {
	return &postgresStorage{db: db}
}

func (s *postgresStorage) Add(t models.Trip, ctx context.Context) error {
	_, span := otel.Tracer(globals.TracerAppName).Start(ctx, "postgres-add-query")
	defer span.End()

	tx := s.db.WithContext(ctx).Create(&t)

	return tx.Error
}

func (s *postgresStorage) Update(t models.Trip, ctx context.Context) error {

	_, span := otel.Tracer(globals.TracerAppName).Start(ctx, "postgres-update-query")
	defer span.End()

	tx := s.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true})

	var points []models.TripPoint
	tx.Where("trip_id = ?", t.ID).Find(&points)

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

	if len(deleted) > 0 {
		deleteTx := tx.Delete(models.NewTripPoint(), deleted)

		if deleteTx.Error != nil {
			return deleteTx.Error
		}
	}

	saveTx := tx.Save(&t)

	return saveTx.Error
}

func (s *postgresStorage) Delete(id int, ctx context.Context) error {
	_, span := otel.Tracer(globals.TracerAppName).Start(ctx, "postgres-delete-query")
	defer span.End()

	tx := s.db.WithContext(ctx).Delete(models.NewTrip(), id)

	return tx.Error
}

func (s *postgresStorage) Get(id int, ctx context.Context) (models.Trip, error) {

	_, span := otel.Tracer(globals.TracerAppName).Start(ctx, "postgres-get-by-id-query")
	defer span.End()

	t := models.NewTrip()

	tx := s.db.WithContext(ctx).Preload("Points").First(t, id)

	return *t, tx.Error
}

func (s *postgresStorage) GetAll(sp *models.TripSearchParam, ctx context.Context) (*[]models.Trip, error) {

	_, span := otel.Tracer(globals.TracerAppName).Start(ctx, "postgrs-get-all-query")
	defer span.End()

	trips := make([]models.Trip, 1)

	o := (sp.Page - 1) * sp.PageSize
	l := sp.PageSize

	tx := s.db.WithContext(ctx).Preload("Points").Offset(o).Limit(l).Find(&trips)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &trips, nil
}
