package repository

import (
	"context"

	"github.com/backend-delery/api/internal/models"
	"gorm.io/gorm"
)

type orderPostgresRepo struct {
	db *gorm.DB
}

// NewOrderRepository creates a new instance of OrderRepository using PostgreSQL.
func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderPostgresRepo{db: db}
}

func (r *orderPostgresRepo) Create(ctx context.Context, order *models.Order) error {
	// Use WithContext to propagate cancellations or timeouts avoiding memory leaks
	return r.db.WithContext(ctx).Create(order).Error
}

func (r *orderPostgresRepo) GetByID(ctx context.Context, id uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.WithContext(ctx).Preload("Store").Preload("Driver").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderPostgresRepo) UpdateStatus(ctx context.Context, id uint, status models.OrderStatus) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderPostgresRepo) AssignDriver(ctx context.Context, orderID uint, driverID uint) error {
	return r.db.WithContext(ctx).Model(&models.Order{}).Where("id = ?", orderID).Updates(map[string]interface{}{
		"driver_id": driverID,
		"status":    models.OrderStatusPreparing, // Implicit business rule here or handle via service. (Keep it simple here)
	}).Error
}

// GetNearbyPending finds orders in PENDING status within a certain radius (in kilometers) from lat/lon.
// It resolves the Haversine formula strictly mathematically over the relationship store.
func (r *orderPostgresRepo) GetNearbyPending(ctx context.Context, lat, lon float64, radiusKm float64) ([]*models.Order, error) {
	var orders []*models.Order

	// Using Haversine formula in Postgres. 6371 is the Earth's radius in KM.
	query := `
	SELECT orders.* 
	FROM orders
	JOIN stores ON stores.id = orders.store_id
	WHERE orders.status = ?
	AND (
		6371 * acos(
			cos(radians(?)) * cos(radians(stores.latitude)) * cos(radians(stores.longitude) - radians(?)) +
			sin(radians(?)) * sin(radians(stores.latitude))
		)
	) <= ?
	`

	err := r.db.WithContext(ctx).Raw(query, models.OrderStatusPending, lat, lon, lat, radiusKm).
		Scan(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}
