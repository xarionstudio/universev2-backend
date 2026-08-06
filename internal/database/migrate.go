package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB, migrationsDir string) error {
	// Create migration tracking table if it doesn't exist
	// This ensures migrations only run once, even across container restarts.
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(50) PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`).Error; err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	// Load already-applied migration versions
	var appliedVersions []string
	if err := db.Table("schema_migrations").Pluck("version", &appliedVersions).Error; err != nil {
		return fmt.Errorf("failed to read applied migrations: %w", err)
	}
	appliedSet := make(map[string]bool)
	for _, v := range appliedVersions {
		appliedSet[v] = true
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	var upFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			upFiles = append(upFiles, entry.Name())
		}
	}

	sort.Strings(upFiles)

	for _, file := range upFiles {
		// Extract version prefix (e.g. "000001" from "000001_init_schema.up.sql")
		version := strings.SplitN(file, "_", 2)[0]

		if appliedSet[version] {
			log.Printf("[Migration] Skipping %s (already applied)", file)
			continue
		}

		fullPath := filepath.Join(migrationsDir, file)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		log.Printf("[Migration] Executing %s ...", file)
		if err := db.Exec(string(content)).Error; err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		// Record migration as applied
		if err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version).Error; err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}

		log.Printf("[Migration] Successfully applied %s", file)
	}

	return nil
}
