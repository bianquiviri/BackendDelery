package main

import (
	"log"

	"github.com/backend-delery/api/database"
	"github.com/backend-delery/api/internal/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// Set Gin mode (Release or Debug)
	gin.SetMode(cfg.GinMode)

	// 2. Initialize Database Connection
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Fatal error initializing DB: %v", err)
	}

	// Safely close the database connection when the application shuts down
	sqlDB, err := db.DB()
	if err == nil {
		defer sqlDB.Close()
	}

	// 2.1 Run Migrations
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Fatal error during DB migration: %v", err)
	}

	// 3. Initialize HTTP Router
	router := gin.Default()

	// Simple health check endpoint
	router.GET("/health", func(c *gin.Context) {
		// Ping DB
		if err := sqlDB.Ping(); err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "Database disconnected"})
			return
		}
		c.JSON(200, gin.H{"status": "ok", "message": "DaaS engine running smoothly"})
	})

	// 4. Start Server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Fatal error starting server: %v", err)
	}
}
