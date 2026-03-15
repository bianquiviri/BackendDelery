package repository

import (
	"context"

	"github.com/backend-delery/api/internal/models"
	"gorm.io/gorm"
)

type storePostgresRepo struct {
	db *gorm.DB
}

// NewStoreRepository creates a new instance of StoreRepository using PostgreSQL.
func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storePostgresRepo{db: db}
}

func (r *storePostgresRepo) GetByID(ctx context.Context, id uint) (*models.Store, error) {
	var store models.Store
	if err := r.db.WithContext(ctx).First(&store, id).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

func (r *storePostgresRepo) Create(ctx context.Context, store *models.Store) error {
	return r.db.WithContext(ctx).Create(store).Error
}
