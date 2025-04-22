package main

import (
	"fmt"
	"os"

	"github.com/arifsetyawan/validra/src/config"
	"github.com/arifsetyawan/validra/src/pkg/database"
	"github.com/arifsetyawan/validra/src/pkg/logger"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()
	log.Info("Starting database migration...")

	// Load configuration
	cfg := config.Load()
	log.Info("Configuration loaded")

	// Initialize PostgreSQL with GORM
	db, err := database.NewPostgresDB(
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	if err != nil {
		log.Error("Failed to connect to PostgreSQL database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	log.Info("Running database migrations...")
	if err := db.Migrate(); err != nil {
		log.Error("Migration failed: %v", err)
		os.Exit(1)
	}

	log.Info("Database migration completed successfully")
	fmt.Println("All migrations have been applied.")
}
