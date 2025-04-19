package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// SQLiteDB encapsulates the database connection
type SQLiteDB struct {
	DB *sql.DB
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &SQLiteDB{DB: db}, nil
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	return s.DB.Close()
}

// Migrate creates necessary tables if they don't exist
func (s *SQLiteDB) Migrate() error {
	log.Println("Running database migrations...")

	// Create Resources table
	_, err := s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS resources (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			attributes BLOB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create resources table: %w", err)
	}

	// Create Actions table
	_, err = s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS actions (
			id TEXT PRIMARY KEY,
			resource_id TEXT NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			attributes BLOB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create actions table: %w", err)
	}

	// Create Roles table
	_, err = s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS roles (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create roles table: %w", err)
	}

	// Create Users table
	_, err = s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			attributes BLOB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create User Sets table
	_, err = s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS user_sets (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			conditions BLOB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create user_sets table: %w", err)
	}

	// Create Resource Sets table
	_, err = s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS resource_sets (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			conditions BLOB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create resource_sets table: %w", err)
	}

	// Create Permissions table
	_, err = s.DB.Exec(`
		CREATE TABLE IF NOT EXISTS permissions (
			id TEXT PRIMARY KEY,
			role_id TEXT NOT NULL,
			user_id TEXT,
			user_set_id TEXT,
			resource_id TEXT,
			resource_set_id TEXT,
			effect TEXT NOT NULL,
			conditions BLOB,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
			FOREIGN KEY (user_set_id) REFERENCES user_sets (id) ON DELETE CASCADE,
			FOREIGN KEY (resource_id) REFERENCES resources (id) ON DELETE CASCADE,
			FOREIGN KEY (resource_set_id) REFERENCES resource_sets (id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create permissions table: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}
