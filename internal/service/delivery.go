package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/backend-delery/api/internal/models"
	"github.com/backend-delery/api/internal/repository"
)

type orderService struct {
	orderRepo  repository.OrderRepository
	storeRepo  repository.StoreRepository
	driverRepo repository.DriverRepository
}

// NewOrderService injects repositories into the service layer fulfilling the DI pattern.
func NewOrderService(or repository.OrderRepository, sr repository.StoreRepository, dr repository.DriverRepository) OrderService {
	return &orderService{
		orderRepo:  or,
		storeRepo:  sr,
		driverRepo: dr,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, storeID uint, total float64, customerAddress string) (*models.Order, error) {
	// 1. Validate Store exists
	store, err := s.storeRepo.GetByID(ctx, storeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch store: %w", err)
	}

	if store == nil {
		return nil, errors.New("store not found")
	}

	order := &models.Order{
		StoreID:         storeID,
		Status:          models.OrderStatusPending,
		Total:           total,
		CustomerAddress: customerAddress,
	}

	// 2. Persist
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	log.Printf("Order %d created successfully for store %d", order.ID, storeID)
	return order, nil
}

func (s *orderService) GetNearbyOrders(ctx context.Context, driverID uint, lat, lon float64) ([]*models.Order, error) {
	// Optionally validate driver exists and update their last known location async or sync
	// For production, location updates might happen via a separate high-frequency endpoint (e.g. WebSocket/UDP)
	_ = s.driverRepo.UpdateLocation(ctx, driverID, lat, lon)

	// Fetch orders within an arbitrary radius (e.g., 5 km)
	const radiusKm = 5.0
	orders, err := s.orderRepo.GetNearbyPending(ctx, lat, lon, radiusKm)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch nearby pending orders: %w", err)
	}

	return orders, nil
}

func (s *orderService) UpdateOrderStatus(ctx context.Context, orderID uint, status models.OrderStatus) error {
	// A real-world DaaS would implement an FSM (Finite State Machine) here to guarantee valid transitions
	// e.g. Validating you can't go from PENDING directly to DELIVERED without PREPARING and ON_ROUTE.
	return s.orderRepo.UpdateStatus(ctx, orderID, status)
}
