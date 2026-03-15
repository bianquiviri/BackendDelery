package repository

import (
	"context"

	"github.com/backend-delery/api/internal/models"
	"gorm.io/gorm"
)

type driverPostgresRepo struct {
	db *gorm.DB
}

// NewDriverRepository creates a new instance of DriverRepository using PostgreSQL.
func NewDriverRepository(db *gorm.DB) DriverRepository {
	return &driverPostgresRepo{db: db}
}

func (r *driverPostgresRepo) GetByID(ctx context.Context, id uint) (*models.Driver, error) {
	var driver models.Driver
	if err := r.db.WithContext(ctx).First(&driver, id).Error; err != nil {
		return nil, err
	}
	return &driver, nil
}

func (r *driverPostgresRepo) UpdateLocation(ctx context.Context, id uint, lat, lon float64) error {
	return r.db.WithContext(ctx).Model(&models.Driver{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_latitude":  lat,
		"last_longitude": lon,
	}).Error
}

func (r *driverPostgresRepo) UpdateStatus(ctx context.Context, id uint, status models.DriverStatus) error {
	return r.db.WithContext(ctx).Model(&models.Driver{}).Where("id = ?", id).Update("status", status).Error
}
