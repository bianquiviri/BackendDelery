package service

import (
	"context"

	"github.com/backend-delery/api/internal/models"
)

// OrderService defines the business logic of Delivery handling.
type OrderService interface {
	CreateOrder(ctx context.Context, storeID uint, total float64, customerAddress string) (*models.Order, error)
	GetNearbyOrders(ctx context.Context, driverID uint, lat, lon float64) ([]*models.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uint, status models.OrderStatus) error
}
