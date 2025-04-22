package atlasmigrate

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
	"github.com/arifsetyawan/validra/src/pkg/logger"
)

// MigrationRunner handles database migrations with Atlas
type MigrationRunner struct {
	db     *sql.DB
	dir    string
	logger *logger.Logger
}

// NewMigrationRunner creates a new MigrationRunner instance
func NewMigrationRunner(db *sql.DB, migrationsDir string, logger *logger.Logger) *MigrationRunner {
	return &MigrationRunner{
		db:     db,
		dir:    migrationsDir,
		logger: logger,
	}
}

// MigrateUp applies all pending migrations
func (m *MigrationRunner) MigrateUp(ctx context.Context) error {
	// Ensure migrations directory exists
	if err := os.MkdirAll(m.dir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Create a postgres driver
	driver, err := postgres.Open(m.db)
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}
	defer driver.Close()

	// Create a directory source for migration files
	dir, err := migrate.NewLocalDir(m.dir)
	if err != nil {
		return fmt.Errorf("failed to create migration directory source: %w", err)
	}

	// Collect all migration files
	files, err := dir.Files()
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}
	if len(files) == 0 {
		m.logger.Info("No migration files found")
		return nil
	}

	// Sort migration files by version
	migrate.SortFiles(files)

	// Create a migration planner
	ex, err := migrate.NewExecutor(driver, dir)
	if err != nil {
		return fmt.Errorf("failed to create migration executor: %w", err)
	}

	// Apply migrations
	m.logger.Info("Applying migrations...")
	if err := ex.ExecuteN(ctx, len(files)); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	m.logger.Info("Migrations applied successfully")
	return nil
}

// CheckMigrationStatus checks the status of all migrations
func (m *MigrationRunner) CheckMigrationStatus(ctx context.Context) error {
	// Create a postgres driver
	driver, err := postgres.Open(m.db)
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}
	defer driver.Close()

	// Create a directory source for migration files
	dir, err := migrate.NewLocalDir(m.dir)
	if err != nil {
		return fmt.Errorf("failed to create migration directory source: %w", err)
	}

	// Get all migration files
	files, err := dir.Files()
	if err != nil {
		return fmt.Errorf("failed to read migration files: %w", err)
	}

	// If there are no migration files, there's nothing to check
	if len(files) == 0 {
		m.logger.Info("No migration files found")
		return nil
	}

	// Sort migration files by version
	migrate.SortFiles(files)

	// Check migration status
	m.logger.Info("Checking migration status...")
	applied, err := driver.MigrationsApplied(ctx)
	if err != nil {
		if errors.Is(err, schema.NotExistError{Name: "atlas_schema_revisions"}) {
			m.logger.Info("No migrations have been applied yet")
			return nil
		}
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Print migration status
	for _, f := range files {
		isApplied := false
		for _, a := range applied {
			if a == f.Version() {
				isApplied = true
				break
			}
		}

		status := "Pending"
		if isApplied {
			status = "Applied"
		}

		m.logger.Info("%s: %s (%s)", f.Version(), filepath.Base(f.Name()), status)
	}

	return nil
}

// GenerateMigration creates a new migration file
func (m *MigrationRunner) GenerateMigration(ctx context.Context, name string) error {
	// TODO: Implement a programmatic way to generate migrations
	// For now, use the CLI tool via the Makefile
	m.logger.Info("Migration generation needs to be done using the atlas CLI tool")
	m.logger.Info("Run: make atlas-migrate-diff")
	return nil
} 