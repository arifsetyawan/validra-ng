package atlasmigrate

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/exec"

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

// MigrateUp applies all pending migrations using the Atlas CLI
func (m *MigrationRunner) MigrateUp(ctx context.Context) error {
	// Ensure migrations directory exists
	if err := os.MkdirAll(m.dir, 0755); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Use Atlas CLI to run migrations (execute atlas migrate apply --env dev)
	m.logger.Info("Applying migrations with Atlas CLI...")
	cmd := exec.CommandContext(ctx, "atlas", "migrate", "apply", "--env", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run Atlas migrations: %w", err)
	}

	m.logger.Info("Migrations applied successfully")
	return nil
}

// CheckMigrationStatus checks the status of all migrations using Atlas CLI
func (m *MigrationRunner) CheckMigrationStatus(ctx context.Context) error {
	// Use Atlas CLI to check migration status (execute atlas migrate status --env dev)
	m.logger.Info("Checking migration status with Atlas CLI...")
	cmd := exec.CommandContext(ctx, "atlas", "migrate", "status", "--env", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	return nil
}

// HashMigrations rehashes migration files to fix checksum errors
func (m *MigrationRunner) HashMigrations(ctx context.Context) error {
	// Use Atlas CLI to rehash migrations (execute atlas migrate hash --env dev)
	m.logger.Info("Rehashing migration files with Atlas CLI...")
	cmd := exec.CommandContext(ctx, "atlas", "migrate", "hash", "--env", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to rehash migration files: %w", err)
	}

	m.logger.Info("Migration files rehashed successfully")
	return nil
}

// GenerateMigration creates a new migration file using Atlas CLI
func (m *MigrationRunner) GenerateMigration(ctx context.Context, name string) error {
	// Use Atlas CLI to generate a migration (execute atlas migrate diff --env dev name)
	m.logger.Info("Generating migration with Atlas CLI...")
	cmd := exec.CommandContext(ctx, "atlas", "migrate", "diff", "--env", "dev", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate migration: %w", err)
	}

	m.logger.Info("Migration generated successfully")
	return nil
} 