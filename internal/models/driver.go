package models

import "time"

type DriverStatus string

const (
	DriverStatusAvailable DriverStatus = "AVAILABLE"
	DriverStatusBusy      DriverStatus = "BUSY"
)

type Driver struct {
	ID            uint         `gorm:"primaryKey" json:"id"`
	Name          string       `gorm:"size:255;not null" json:"name"`
	Vehicle       string       `gorm:"size:100;not null" json:"vehicle"`
	Status        DriverStatus `gorm:"type:varchar(20);default:'AVAILABLE'" json:"status"`
	LastLatitude  float64      `gorm:"not null" json:"last_latitude"`
	LastLongitude float64      `gorm:"not null" json:"last_longitude"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}
