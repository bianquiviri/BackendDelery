package database

import (
	"log"

	"github.com/backend-delery/api/internal/models"
	"gorm.io/gorm"
)

// AutoMigrate runs GORM auto-migrations for all core domain models.
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running database migrations...")
	return db.AutoMigrate(
		&models.Store{},
		&models.Driver{},
		&models.Order{},
	)
}
