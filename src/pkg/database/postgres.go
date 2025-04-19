package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresDB encapsulates the database connection
type PostgresDB struct {
	DB *gorm.DB
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(host, user, password, dbname string, port int, sslmode string) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		host, user, password, dbname, port, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &PostgresDB{DB: db}, nil
}

// Close closes the database connection
func (p *PostgresDB) Close() error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database: %w", err)
	}
	return sqlDB.Close()
}

// Migrate creates necessary tables if they don't exist
func (p *PostgresDB) Migrate() error {
	log.Println("Running database migrations...")

	// Define models for migration
	type Resource struct {
		ID          string `gorm:"primaryKey"`
		Name        string `gorm:"not null"`
		Description string
		Attributes  []byte
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	type Role struct {
		ID          string `gorm:"primaryKey"`
		Name        string `gorm:"not null"`
		Description string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}

	type User struct {
		ID        string `gorm:"primaryKey"`
		Username  string `gorm:"not null;uniqueIndex"`
		Email     string `gorm:"uniqueIndex"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	// Run migrations
	err := p.DB.AutoMigrate(&Resource{}, &Role{}, &User{})
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
