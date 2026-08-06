package main

import (
	"log"
	"os"
	"path/filepath"

	"universev/internal/config"
	"universev/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewGormDB(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Find migrations directory
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		cwd, _ := os.Getwd()
		migrationsDir = filepath.Join(cwd, "migrations")
		if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
			migrationsDir = "../migrations"
		}
	}

	log.Printf("[Migrate] Running migrations from %s ...", migrationsDir)
	if err := database.RunMigrations(db, migrationsDir); err != nil {
		log.Fatalf("[Migrate] Migration failed: %v", err)
	}

	log.Println("[Migrate] All migrations and seeds executed successfully!")
}
