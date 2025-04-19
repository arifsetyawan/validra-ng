package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/arifsetyawan/validra/src/internal/domain"
	"github.com/arifsetyawan/validra/src/pkg/database"
	"github.com/google/uuid"
)

// SQLiteRoleRepository implements domain.RoleRepository using SQLite
type SQLiteRoleRepository struct {
	db *database.SQLiteDB
}

// NewSQLiteRoleRepository creates a new SQLite repository for roles
func NewSQLiteRoleRepository(db *database.SQLiteDB) domain.RoleRepository {
	return &SQLiteRoleRepository{
		db: db,
	}
}

// Create inserts a new role into the database
func (r *SQLiteRoleRepository) Create(ctx context.Context, role *domain.Role) error {
	// Generate a new UUID if not provided
	if role.ID == "" {
		role.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	query := `
		INSERT INTO roles (id, name, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.DB.ExecContext(
		ctx,
		query,
		role.ID,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

// GetByID retrieves a role by ID
func (r *SQLiteRoleRepository) GetByID(ctx context.Context, id string) (*domain.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = ?
	`
	row := r.db.DB.QueryRowContext(ctx, query, id)

	var role domain.Role
	err := row.Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &role, nil
}

// List retrieves a paginated list of roles
func (r *SQLiteRoleRepository) List(ctx context.Context, limit, offset int) ([]*domain.Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var roles []*domain.Role
	for rows.Next() {
		var role domain.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, &role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating roles: %w", err)
	}

	return roles, nil
}

// Update updates an existing role
func (r *SQLiteRoleRepository) Update(ctx context.Context, role *domain.Role) error {
	// Update the timestamp
	role.UpdatedAt = time.Now()

	query := `
		UPDATE roles
		SET name = ?, description = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.DB.ExecContext(
		ctx,
		query,
		role.Name,
		role.Description,
		role.UpdatedAt,
		role.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

// Delete deletes a role by ID
func (r *SQLiteRoleRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM roles
		WHERE id = ?
	`
	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}
