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

// SQLiteResourceRepository implements domain.ResourceRepository using SQLite
type SQLiteResourceRepository struct {
	db *database.SQLiteDB
}

// NewSQLiteResourceRepository creates a new SQLite repository for resources
func NewSQLiteResourceRepository(db *database.SQLiteDB) domain.ResourceRepository {
	return &SQLiteResourceRepository{
		db: db,
	}
}

// Create inserts a new resource into the database
func (r *SQLiteResourceRepository) Create(ctx context.Context, resource *domain.Resource) error {
	// Generate a new UUID if ID is not provided
	if resource.ID == "" {
		resource.ID = uuid.New().String()
	}

	now := time.Now()
	resource.CreatedAt = now
	resource.UpdatedAt = now

	query := `
		INSERT INTO resources (id, name, description, attributes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.DB.ExecContext(
		ctx,
		query,
		resource.ID,
		resource.Name,
		resource.Description,
		resource.Attributes,
		resource.CreatedAt,
		resource.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	return nil
}

// GetByID retrieves a resource by ID
func (r *SQLiteResourceRepository) GetByID(ctx context.Context, id string) (*domain.Resource, error) {
	query := `
		SELECT id, name, description, attributes, created_at, updated_at
		FROM resources
		WHERE id = ?
	`

	var resource domain.Resource
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&resource.ID,
		&resource.Name,
		&resource.Description,
		&resource.Attributes,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("resource not found")
		}
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	return &resource, nil
}

// List retrieves a paginated list of resources
func (r *SQLiteResourceRepository) List(ctx context.Context, limit, offset int) ([]*domain.Resource, error) {
	query := `
		SELECT id, name, description, attributes, created_at, updated_at
		FROM resources
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list resources: %w", err)
	}
	defer rows.Close()

	var resources []*domain.Resource
	for rows.Next() {
		var resource domain.Resource
		err := rows.Scan(
			&resource.ID,
			&resource.Name,
			&resource.Description,
			&resource.Attributes,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan resource: %w", err)
		}
		resources = append(resources, &resource)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}

	return resources, nil
}

// Update updates an existing resource
func (r *SQLiteResourceRepository) Update(ctx context.Context, resource *domain.Resource) error {
	resource.UpdatedAt = time.Now()

	query := `
		UPDATE resources
		SET name = ?, description = ?, attributes = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.DB.ExecContext(
		ctx,
		query,
		resource.Name,
		resource.Description,
		resource.Attributes,
		resource.UpdatedAt,
		resource.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update resource: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("resource not found")
	}

	return nil
}

// Delete deletes a resource by ID
func (r *SQLiteResourceRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM resources WHERE id = ?`

	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("resource not found")
	}

	return nil
}
