package models

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPreparing OrderStatus = "PREPARING"
	OrderStatusOnRoute   OrderStatus = "ON_ROUTE"
	OrderStatusDelivered OrderStatus = "DELIVERED"
)

type Order struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	StoreID         uint        `gorm:"not null" json:"store_id"`
	Store           Store       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"store"`  // Belongs To
	DriverID        *uint       `json:"driver_id"`                                                    // Nullable because it is unassigned initially
	Driver          *Driver     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"driver"` // Optional Belongs To
	Status          OrderStatus `gorm:"type:varchar(20);default:'PENDING'" json:"status"`
	Total           float64     `gorm:"type:decimal(10,2);not null" json:"total"`
	CustomerAddress string      `gorm:"not null" json:"customer_address"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}
