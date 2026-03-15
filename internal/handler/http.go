package handler

import (
	"net/http"

	"github.com/backend-delery/api/internal/models"
	"github.com/backend-delery/api/internal/service"
	"github.com/gin-gonic/gin"
)

// OrderHandler handles HTTP requests for orders.
type OrderHandler struct {
	svc service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// RegisterRoutes attaches the handler methods to the Gin router.
func (h *OrderHandler) RegisterRoutes(router *gin.Engine) {
	orders := router.Group("/orders")
	{
		orders.POST("", h.CreateOrder)
		orders.GET("/nearby", h.GetNearbyOrders)
		orders.PATCH("/:id/status", h.UpdateOrderStatus)
	}
}

// CreateOrderRequest defines the expected JSON payload for creating an order.
type CreateOrderRequest struct {
	StoreID         uint    `json:"store_id" binding:"required"`
	Total           float64 `json:"total" binding:"required,gt=0"`
	CustomerAddress string  `json:"customer_address" binding:"required"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// The context from Gin `c.Request.Context()` is passed down.
	// If the client drops the connection, the DB query will be safely canceled avoiding resource leaks.
	order, err := h.svc.CreateOrder(c.Request.Context(), req.StoreID, req.Total, req.CustomerAddress)
	if err != nil {
		// Log the internal error, but return a sanitized message to the client (Security)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, order)
}

// GetNearbyOrdersRequest defines the query parameters for fetching nearby orders.
type GetNearbyOrdersRequest struct {
	DriverID  uint    `form:"driver_id" binding:"required"`
	Latitude  float64 `form:"lat" binding:"required"`
	Longitude float64 `form:"lon" binding:"required"`
}

func (h *OrderHandler) GetNearbyOrders(c *gin.Context) {
	var req GetNearbyOrdersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters", "details": err.Error()})
		return
	}

	orders, err := h.svc.GetNearbyOrders(c.Request.Context(), req.DriverID, req.Latitude, req.Longitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query nearby orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

// UpdateOrderStatusRequest defines the JSON payload for updating an order status.
type UpdateOrderStatusRequest struct {
	Status models.OrderStatus `json:"status" binding:"required"`
}

// UpdateOrderStatus handles the status transition of an order.
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	// Parse Order ID from URL
	var params struct {
		ID uint `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	if err := h.svc.UpdateOrderStatus(c.Request.Context(), params.ID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated successfully"})
}
