package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/nitefawkes/ham-radio-cloud/internal/config"
	"github.com/nitefawkes/ham-radio-cloud/internal/handlers"
	"github.com/nitefawkes/ham-radio-cloud/internal/repositories"
	"github.com/nitefawkes/ham-radio-cloud/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := config.ConnectDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	qsoRepo := repositories.NewQSORepository(db)
	propRepo := repositories.NewPropagationRepository(db)
	sdrRepo := repositories.NewSDRRepository(db)

	// Initialize services
	qsoService := services.NewQSOService(qsoRepo)
	propService := services.NewPropagationService(propRepo)
	sdrService := services.NewSDRService(sdrRepo)

	// Initialize handlers
	qsoHandler := handlers.NewQSOHandler(qsoService)
	adifHandler := handlers.NewADIFHandler(qsoService)
	propHandler := handlers.NewPropagationHandler(propService)
	sdrHandler := handlers.NewSDRHandler(sdrService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Ham Radio Cloud API v1.0",
		ServerHeader: "Ham-Radio-Cloud",
		ErrorHandler: handlers.ErrorHandler,
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check endpoint
	app.Get("/health", handlers.HealthCheck)
	app.Get("/api/v1/health", handlers.HealthCheck)

	// API v1 routes
	v1 := app.Group("/api/v1")
	handlers.RegisterRoutes(v1, qsoHandler, adifHandler, propHandler, sdrHandler)

	// Fetch initial propagation data
	go func() {
		log.Println("Fetching initial propagation data...")
		if err := propService.FetchAndStore(); err != nil {
			log.Printf("Warning: Failed to fetch initial propagation data: %v", err)
		}
	}()

	// Start propagation data refresh scheduler (every 15 minutes)
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Scheduled propagation data refresh...")
			if err := propService.FetchAndStore(); err != nil {
				log.Printf("Warning: Scheduled propagation refresh failed: %v", err)
			}
		}
	}()

	// Cleanup old propagation data daily
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Cleaning up old propagation data...")
			if err := propService.CleanupOldData(30); err != nil { // Keep 30 days
				log.Printf("Warning: Failed to cleanup old data: %v", err)
			}
		}
	}()

	// Fetch initial SDR directory
	go func() {
		log.Println("Fetching initial SDR directory...")
		if err := sdrService.RefreshDirectory(); err != nil {
			log.Printf("Warning: Failed to fetch initial SDR directory: %v", err)
		}
	}()

	// Start SDR directory refresh scheduler (every 6 hours)
	go func() {
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("Scheduled SDR directory refresh...")
			if err := sdrService.RefreshDirectory(); err != nil {
				log.Printf("Warning: Scheduled SDR refresh failed: %v", err)
			}
		}
	}()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Ham Radio Cloud API starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
