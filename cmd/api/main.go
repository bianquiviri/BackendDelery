package main

import (
	"log"

	"os"

	"github.com/backend-delery/api/database"
	"github.com/backend-delery/api/internal/config"
	"github.com/backend-delery/api/internal/handler"
	"github.com/backend-delery/api/internal/repository"
	"github.com/backend-delery/api/internal/service"
	"github.com/gin-gonic/gin"

	_ "github.com/backend-delery/api/docs" // Swagger docs
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           DaaS Backend API
// @version         1.0
// @description     Delivery as a Service Backend API documentation.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8084
// @BasePath  /

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

	// 2.2 Seed Data if requested
	if os.Getenv("SEED") == "true" {
		if err := database.SeedData(db); err != nil {
			log.Printf("Warning: Seeding failed: %v", err)
		}
	}

	// 3. Dependency Injection (DI) Setup
	// 3.1 Repositories
	storeRepo := repository.NewStoreRepository(db)
	driverRepo := repository.NewDriverRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// 3.2 Services
	orderSvc := service.NewOrderService(orderRepo, storeRepo, driverRepo)

	// 3.3 Handlers
	orderHandler := handler.NewOrderHandler(orderSvc)

	// 4. Initialize HTTP Router
	router := gin.Default()

	// 4.1 Global Middleware
	// Recovery middleware recovers from any panics and writes a 500 if there was one. (Resilience)
	router.Use(gin.Recovery())

	// 4.2 Register Routes
	orderHandler.RegisterRoutes(router)

	// 4.3 Serve Static Files (Documentation and API Explorer)
	router.StaticFile("/", "./public/index.html")
	router.Static("/assets", "./public/assets")
	router.StaticFS("/public", gin.Dir("public", false))

	// 4.4 Swagger Documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
