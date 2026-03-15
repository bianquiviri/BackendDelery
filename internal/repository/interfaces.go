package repository

import (
	"context"

	"github.com/backend-delery/api/internal/models"
)

// StoreRepository defines the persistence interface for Store entity.
type StoreRepository interface {
	GetByID(ctx context.Context, id uint) (*models.Store, error)
	Create(ctx context.Context, store *models.Store) error
}

// DriverRepository defines the persistence interface for Driver entity.
type DriverRepository interface {
	GetByID(ctx context.Context, id uint) (*models.Driver, error)
	UpdateLocation(ctx context.Context, id uint, lat, lon float64) error
	UpdateStatus(ctx context.Context, id uint, status models.DriverStatus) error
}

// OrderRepository defines the persistence interface for Order entity.
type OrderRepository interface {
	Create(ctx context.Context, order *models.Order) error
	GetByID(ctx context.Context, id uint) (*models.Order, error)
	GetNearbyPending(ctx context.Context, lat, lon float64, radiusKm float64) ([]*models.Order, error)
	UpdateStatus(ctx context.Context, id uint, status models.OrderStatus) error
	AssignDriver(ctx context.Context, orderID uint, driverID uint) error
}
