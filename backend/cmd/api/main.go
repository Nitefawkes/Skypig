package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/Nitefawkes/Skypig/backend/internal/config"
	"github.com/Nitefawkes/Skypig/backend/internal/database"
	"github.com/Nitefawkes/Skypig/backend/internal/handlers"
	"github.com/Nitefawkes/Skypig/backend/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := database.NewUserRepository(db)
	qsoRepo := database.NewQSORepository(db)

	// Initialize services
	qsoService := services.NewQSOService(qsoRepo, userRepo)
	adifService := services.NewADIFService(qsoRepo, userRepo, qsoService)
	propagationService := services.NewPropagationService(true) // Use mock data for now

	// Initialize handlers
	qsoHandler := handlers.NewQSOHandler(qsoService)
	adifHandler := handlers.NewADIFHandler(adifService)
	propagationHandler := handlers.NewPropagationHandler(propagationService)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Ham-Radio Cloud API",
		ServerHeader: "Ham-Radio-Cloud",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Server.CORSOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		// Check database health
		if err := db.HealthCheck(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "error",
				"service": "ham-radio-cloud-api",
				"error":   "database unhealthy",
			})
		}

		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "ham-radio-cloud-api",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := app.Group("/api/v1")

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Ham-Radio Cloud API v1",
			"endpoints": []string{
				"/health",
				"/api/v1/qsos",
				"/api/v1/qsos/:id",
				"/api/v1/qsos/stats",
				"/api/v1/qsos/import",
				"/api/v1/qsos/export",
				"/api/v1/propagation",
				"/api/v1/propagation/forecast",
				"/api/v1/sdr (coming soon)",
			},
		})
	})

	// QSO routes
	qsos := v1.Group("/qsos")
	qsos.Get("/stats", qsoHandler.GetStats)
	qsos.Post("/import", adifHandler.ImportADIF)
	qsos.Get("/export", adifHandler.ExportADIF)
	qsos.Post("/validate", adifHandler.ValidateADIF)
	qsos.Get("/", qsoHandler.ListQSOs)
	qsos.Post("/", qsoHandler.CreateQSO)
	qsos.Get("/:id", qsoHandler.GetQSO)
	qsos.Put("/:id", qsoHandler.UpdateQSO)
	qsos.Delete("/:id", qsoHandler.DeleteQSO)

	// Propagation routes
	propagation := v1.Group("/propagation")
	propagation.Get("/", propagationHandler.GetCurrent)
	propagation.Get("/forecast", propagationHandler.GetForecast)

	// Start server
	log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
	log.Printf("📊 Environment: %s", cfg.Server.Environment)
	log.Printf("🗄️  Database: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	if err := app.Listen(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
