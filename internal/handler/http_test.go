package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/backend-delery/api/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockOrderService is a mock implementation of service.OrderService
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, storeID uint, total float64, address string) (*models.Order, error) {
	args := m.Called(ctx, storeID, total, address)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) GetNearbyOrders(ctx context.Context, driverID uint, lat, lon float64) ([]*models.Order, error) {
	args := m.Called(ctx, driverID, lat, lon)
	return args.Get(0).([]*models.Order), args.Error(1)
}

func (m *MockOrderService) UpdateOrderStatus(ctx context.Context, id uint, status models.OrderStatus) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestOrderHandler_CreateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockOrderService)
		h := NewOrderHandler(mockSvc)
		
		expectedOrder := &models.Order{ID: 1, StoreID: 10, Total: 25.5, CustomerAddress: "123 Test St"}
		mockSvc.On("CreateOrder", mock.Anything, uint(10), 25.5, "123 Test St").Return(expectedOrder, nil)
		
		router := gin.New()
		h.RegisterRoutes(router)
		
		body, _ := json.Marshal(CreateOrderRequest{
			StoreID:         10,
			Total:           25.5,
			CustomerAddress: "123 Test St",
		})
		
		req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusCreated, w.Code)
		
		var resp models.Order
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, expectedOrder.ID, resp.ID)
		mockSvc.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		h := NewOrderHandler(nil)
		router := gin.New()
		h.RegisterRoutes(router)
		
		req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString("invalid json"))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestOrderHandler_GetNearbyOrders(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockOrderService)
		h := NewOrderHandler(mockSvc)
		
		expectedOrders := []*models.Order{{ID: 1}, {ID: 2}}
		mockSvc.On("GetNearbyOrders", mock.Anything, uint(1), 40.0, -70.0).Return(expectedOrders, nil)
		
		router := gin.New()
		h.RegisterRoutes(router)
		
		req, _ := http.NewRequest(http.MethodGet, "/orders/nearby?driver_id=1&lat=40.0&lon=-70.0", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestOrderHandler_UpdateOrderStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockOrderService)
		h := NewOrderHandler(mockSvc)
		
		mockSvc.On("UpdateOrderStatus", mock.Anything, uint(1), models.OrderStatusDelivered).Return(nil)
		
		router := gin.New()
		h.RegisterRoutes(router)
		
		body, _ := json.Marshal(UpdateOrderStatusRequest{Status: models.OrderStatusDelivered})
		req, _ := http.NewRequest(http.MethodPatch, "/orders/1/status", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
