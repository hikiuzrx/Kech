package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/smartwaste/backend/internal/config"
	"github.com/smartwaste/backend/internal/database"
	"github.com/smartwaste/backend/internal/handlers"
	"github.com/smartwaste/backend/internal/mqtt"
	"github.com/smartwaste/backend/internal/repository"
	"github.com/smartwaste/backend/internal/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize database connection
	db, err := database.InitDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	driverRepo := repository.NewDriverRepository(db)
	binRepo := repository.NewBinRepository(db)
	collectionRepo := repository.NewCollectionRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	pricingRepo := repository.NewPricingRepository(db)

	// Initialize services
	notificationSvc := services.NewNotificationService(driverRepo)
	valuationSvc := services.NewValuationService(pricingRepo)
	routeSvc := services.NewRouteService(binRepo, &cfg.Google)
	analyticsSvc := services.NewAnalyticsService(binRepo, collectionRepo, driverRepo)

	// Initialize MQTT client
	mqttClient := mqtt.NewClient(&cfg.MQTT, binRepo, notificationSvc)
	if err := mqttClient.Connect(); err != nil {
		log.Printf("Warning: Failed to connect to MQTT broker: %v", err)
		log.Println("Continuing without MQTT - IoT data ingestion will be unavailable")
	} else {
		defer mqttClient.Disconnect()
		if err := mqttClient.Subscribe(); err != nil {
			log.Printf("Warning: Failed to subscribe to MQTT topics: %v", err)
		}
	}

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo)
	driverHandler := handlers.NewDriverHandler(driverRepo, binRepo, collectionRepo, routeSvc)
	binHandler := handlers.NewBinHandler(binRepo)
	companyHandler := handlers.NewCompanyHandler(companyRepo, pricingRepo, valuationSvc)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsSvc)

	// Setup router
	router := setupRouter(userHandler, driverHandler, binHandler, companyHandler, analyticsHandler, mqttClient)

	// Create server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func setupRouter(
	userHandler *handlers.UserHandler,
	driverHandler *handlers.DriverHandler,
	binHandler *handlers.BinHandler,
	companyHandler *handlers.CompanyHandler,
	analyticsHandler *handlers.AnalyticsHandler,
	mqttClient *mqtt.Client,
) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(handlers.RecoveryMiddleware())
	router.Use(handlers.LoggerMiddleware())
	router.Use(handlers.CORSMiddleware())
	router.Use(handlers.RequestIDMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		status := "healthy"
		mqttStatus := "connected"
		if mqttClient == nil || !mqttClient.IsConnected() {
			mqttStatus = "disconnected"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":      status,
			"mqtt_status": mqttStatus,
			"timestamp":   time.Now().UTC(),
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.GET("", userHandler.ListUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("/:id/rewards", userHandler.GetRewardPoints)
			users.POST("/:id/rewards", userHandler.AddRewardPoints)
		}

		// Driver routes
		drivers := v1.Group("/drivers")
		{
			drivers.GET("", driverHandler.ListDrivers)
			drivers.POST("", driverHandler.CreateDriver)
			drivers.GET("/:id", driverHandler.GetDriver)
			drivers.PUT("/:id", driverHandler.UpdateDriver)
			drivers.PUT("/:id/location", driverHandler.UpdateLocation)
			drivers.GET("/:id/routes", driverHandler.GetRoutes)
			drivers.POST("/:id/verify", driverHandler.VerifyTask)
			drivers.GET("/:id/stats", driverHandler.GetDriverStats)
		}

		// Bin routes
		bins := v1.Group("/bins")
		{
			bins.GET("", binHandler.ListBins)
			bins.POST("", binHandler.CreateBin)
			bins.GET("/needs-collection", binHandler.GetBinsNeedingCollection)
			bins.GET("/statistics", binHandler.GetBinStatistics)
			bins.GET("/:id", binHandler.GetBin)
			bins.PUT("/:id", binHandler.UpdateBin)
			bins.DELETE("/:id", binHandler.DeleteBin)
		}

		// Company routes
		companies := v1.Group("/companies")
		{
			companies.GET("", companyHandler.ListCompanies)
			companies.POST("", companyHandler.CreateCompany)
			companies.GET("/:id", companyHandler.GetCompany)
			companies.PUT("/:id", companyHandler.UpdateCompany)
			companies.DELETE("/:id", companyHandler.DeleteCompany)
		}

		// Pricing rules routes
		pricingRules := v1.Group("/pricing-rules")
		{
			pricingRules.GET("", companyHandler.ListPricingRules)
			pricingRules.POST("", companyHandler.CreatePricingRule)
			pricingRules.GET("/:id", companyHandler.GetPricingRule)
			pricingRules.PUT("/:id", companyHandler.UpdatePricingRule)
			pricingRules.DELETE("/:id", companyHandler.DeletePricingRule)
		}

		// Valuations
		v1.POST("/valuations", companyHandler.CalculateValuation)

		// Analytics routes
		analytics := v1.Group("/analytics")
		{
			analytics.GET("/dashboard", analyticsHandler.GetDashboardStats)
			analytics.GET("/bins", analyticsHandler.GetBinAnalytics)
			analytics.GET("/drivers", analyticsHandler.GetDriverAnalytics)
			analytics.GET("/collections", analyticsHandler.GetCollectionAnalytics)
		}
	}

	return router
}
