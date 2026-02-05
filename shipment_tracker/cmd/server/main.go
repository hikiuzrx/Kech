package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/smartwaste/shipment-tracker/internal/config"
	"github.com/smartwaste/shipment-tracker/internal/database"
	"github.com/smartwaste/shipment-tracker/internal/handlers"
	"github.com/smartwaste/shipment-tracker/internal/nats"
	"github.com/smartwaste/shipment-tracker/internal/repository"
	"github.com/smartwaste/shipment-tracker/internal/services"
)

func main() {
	// 1. Load Configuration
	cfg := config.LoadConfig()

	// 2. Initialize Database
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// 3. Initialize NATS
	natsClient := nats.NewClient(&cfg.NATS)
	if err := natsClient.Connect(); err != nil {
		log.Printf("Warning: Failed to connect to NATS: %v. Continuing without messaging...", err)
	} else {
		defer natsClient.Close()
	}

	// 4. Initialize Repositories
	shipmentRepo := repository.NewShipmentRepository(db)
	transitionRepo := repository.NewTransitionRepository(db)
	// contractRepo := repository.NewContractRepository(db) // For later

	// 5. Initialize Services
	shipmentService := services.NewShipmentService(shipmentRepo, transitionRepo, natsClient)

	// 6. Initialize Handlers
	shipmentHandler := handlers.NewShipmentHandler(shipmentService)

	// 7. Setup Router
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		shipments := v1.Group("/shipments")
		{
			shipments.POST("", shipmentHandler.CreateShipment)
			shipments.GET("/:id", shipmentHandler.GetShipment)
			shipments.POST("/:id/assign-driver", shipmentHandler.AssignDriver)
		}
	}

	// 8. Start Server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
