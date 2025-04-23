package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/arifsetyawan/validra/src/config"
	"github.com/arifsetyawan/validra/src/pkg/atlasmigrate"
	"github.com/arifsetyawan/validra/src/pkg/logger"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	// Command line flags
	checkStatus := flag.Bool("status", false, "Check migration status")
	fixChecksums := flag.Bool("fix-checksums", false, "Fix migration checksums")
	flag.Parse()

	// Initialize logger
	log := logger.NewLogger()
	log.Info("Starting Validra migrations...")

	// Load configuration
	cfg := config.Load()
	log.Info("Configuration loaded")

	// Create DB connection string
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Open a direct database connection (not using GORM)
	db, err := openDatabase(dsn)
	if err != nil {
		log.Error("Failed to connect to PostgreSQL database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// Set up the migrations directory path
	migrationsDir := filepath.Join("migrations", "migrations")
	
	// Create a migration runner
	runner := atlasmigrate.NewMigrationRunner(db, migrationsDir, log)
	
	ctx := context.Background()

	// Check status, fix checksums, or apply migrations based on flags
	if *checkStatus {
		log.Info("Checking migration status...")
		if err := runner.CheckMigrationStatus(ctx); err != nil {
			log.Error("Failed to check migration status: %v", err)
			os.Exit(1)
		}
	} else if *fixChecksums {
		log.Info("Fixing migration checksums...")
		if err := runner.HashMigrations(ctx); err != nil {
			log.Error("Failed to fix migration checksums: %v", err)
			os.Exit(1)
		}
		log.Info("Migration checksums fixed successfully")
	} else {
		log.Info("Applying migrations...")
		if err := runner.MigrateUp(ctx); err != nil {
			log.Error("Failed to apply migrations: %v", err)
			os.Exit(1)
		}
		log.Info("Migrations completed successfully")
	}
}

// openDatabase opens a direct database connection
func openDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	
	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return db, nil
}
