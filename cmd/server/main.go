package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"

	"universev/internal/config"
	"universev/internal/database"
	"universev/internal/repository"
	"universev/internal/router"
	"universev/internal/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to GORM PostgreSQL database
	db, err := database.NewGormDB(cfg)
	if err != nil {
		log.Printf("[Warning] Database connection failed: %v.", err)
	}

	// Run migrations automatically
	if db != nil {
		migrationsDir := os.Getenv("MIGRATIONS_DIR")
		if migrationsDir == "" {
			cwd, _ := os.Getwd()
			migrationsDir = filepath.Join(cwd, "migrations")
			if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
				migrationsDir = "../migrations"
			}
		}
		log.Printf("[Server] Running migrations from %s ...", migrationsDir)
		if err := database.RunMigrations(db, migrationsDir); err != nil {
			log.Printf("[Warning] Migration failed: %v", err)
		}
	}

	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CORSAllowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": cfg.AppName,
			"env":     cfg.AppEnv,
		})
	})

	// Initialize Fingerprint Worker
	var fpWorker *worker.FingerprintWorker
	if db != nil {
		fpRepo := repository.NewFingerprintRepo(db)
		attRepo := repository.NewAttendanceRepo(db)
		fpWorker = worker.NewFingerprintWorker(fpRepo, attRepo)

		if cfg.FingerprintEnabled {
			fpWorker.Start(60 * time.Second)
			defer fpWorker.Stop()
			log.Println("[Server] Fingerprint worker started (polling every 60s).")
		} else {
			log.Println("[Server] Fingerprint worker disabled (FINGERPRINT_ENABLED=false).")
		}
	}

	// Setup all API v1 routes with GORM DB
	router.SetupRoutes(app, cfg, db, fpWorker)

	log.Printf("Server starting on port %s", cfg.AppPort)
	if err := app.Listen(":" + cfg.AppPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
