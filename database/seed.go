package database

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/backend-delery/api/internal/models"
	"gorm.io/gorm"
)

// SeedData populates the database with initial test data.
func SeedData(db *gorm.DB) error {
	log.Println("🌱 Seeding database with test data...")

	// 1. Create Test Stores
	stores := []models.Store{
		{Name: "Central Pizza", Address: "Main St 101", Latitude: 40.7128, Longitude: -74.0060, APIKey: "pizza-key-123"},
		{Name: "Burger Haven", Address: "Broadway 50", Latitude: 40.7589, Longitude: -73.9851, APIKey: "burger-key-456"},
		{Name: "Sushi Express", Address: "5th Ave 200", Latitude: 40.7850, Longitude: -73.9683, APIKey: "sushi-key-789"},
	}

	for i := range stores {
		if err := db.FirstOrCreate(&stores[i], models.Store{APIKey: stores[i].APIKey}).Error; err != nil {
			return fmt.Errorf("failed to seed store: %v", err)
		}
	}

	// 2. Create Test Drivers
	drivers := []models.Driver{
		{Name: "John Doe", Vehicle: "Bicycle", Status: models.DriverStatusAvailable, LastLatitude: 40.7130, LastLongitude: -74.0070},
		{Name: "Jane Smith", Vehicle: "Motorcycle", Status: models.DriverStatusAvailable, LastLatitude: 40.7590, LastLongitude: -73.9860},
		{Name: "Bob Wilson", Vehicle: "Car", Status: models.DriverStatusAvailable, LastLatitude: 40.7860, LastLongitude: -73.9690},
	}

	for i := range drivers {
		if err := db.FirstOrCreate(&drivers[i], models.Driver{Name: drivers[i].Name}).Error; err != nil {
			return fmt.Errorf("failed to seed driver: %v", err)
		}
	}

	// 3. Create Bulk Random Orders (if none exist or just add more)
	var count int64
	db.Model(&models.Order{}).Count(&count)
	if count < 50 {
		log.Printf("📦 Creating 100 random orders for testing...")
		for i := 0; i < 100; i++ {
			store := stores[rand.Intn(len(stores))]
			order := models.Order{
				StoreID:         store.ID,
				Total:           float64(rand.Intn(10000)) / 100.0,
				Status:          models.OrderStatusPending,
				CustomerAddress: fmt.Sprintf("Test Street %d", rand.Intn(1000)),
				CreatedAt:       time.Now().Add(time.Duration(-rand.Intn(48)) * time.Hour), // Over last 2 days
			}
			if err := db.Create(&order).Error; err != nil {
				log.Printf("Warning: failed to create seed order: %v", err)
			}
		}
	}

	log.Println("✅ Database seeding completed successfully.")
	return nil
}
