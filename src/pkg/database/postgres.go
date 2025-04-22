package database

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/arifsetyawan/validra/src/internal/model"
	"github.com/arifsetyawan/validra/src/pkg/atlasmigrate"
	"github.com/arifsetyawan/validra/src/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// PostgresDB encapsulates the database connection
type PostgresDB struct {
	DB     *gorm.DB
	logger *logger.Logger
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(host, user, password, dbname string, port int, sslmode string, log *logger.Logger) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(25)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &PostgresDB{
		DB:     db,
		logger: log,
	}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}
	return sqlDB.Close()
}

// Migrate runs database migrations using Atlas
func (p *PostgresDB) Migrate() error {
	p.logger.Info("Running database migrations with Atlas...")

	// Get the SQL connection
	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	// Set up the migrations directory path
	migrationsDir := filepath.Join("migrations", "migrations")
	
	// Create a migration runner
	runner := atlasmigrate.NewMigrationRunner(sqlDB, migrationsDir, p.logger)
	
	// Apply migrations
	ctx := context.Background()
	if err := runner.MigrateUp(ctx); err != nil {
		return fmt.Errorf("failed to run Atlas migrations: %w", err)
	}

	p.logger.Info("PostgreSQL database migrations completed")
	return nil
}

// CheckMigrationStatus checks the status of all migrations
func (p *PostgresDB) CheckMigrationStatus() error {
	p.logger.Info("Checking migration status...")

	// Get the SQL connection
	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}

	// Set up the migrations directory path
	migrationsDir := filepath.Join("migrations", "migrations")
	
	// Create a migration runner
	runner := atlasmigrate.NewMigrationRunner(sqlDB, migrationsDir, p.logger)
	
	// Check migration status
	ctx := context.Background()
	if err := runner.CheckMigrationStatus(ctx); err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	return nil
}

// RunGORMMigrations creates necessary tables using GORM AutoMigrate (legacy method)
// This will be deprecated in favor of Atlas migrations
func (p *PostgresDB) RunGORMMigrations() error {
	p.logger.Info("Running GORM migrations (deprecated)...")

	// Run migrations using models from the model package
	err := p.DB.AutoMigrate(
		&model.Resource{},
		&model.ResourceActions{},
		&model.ResourceABACOptions{},
		&model.ResourceReBACOptions{},
		&model.ResourceSet{},
		&model.Role{},
		&model.User{},
		&model.UserSet{},
		&model.Permission{},
	)
	if err != nil {
		return fmt.Errorf("failed to run GORM migrations: %w", err)
	}

	return nil
}
