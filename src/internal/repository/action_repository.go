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

// SQLiteActionRepository implements domain.ActionRepository using SQLite
type SQLiteActionRepository struct {
	db *database.SQLiteDB
}

// NewSQLiteActionRepository creates a new SQLite repository for actions
func NewSQLiteActionRepository(db *database.SQLiteDB) domain.ActionRepository {
	return &SQLiteActionRepository{
		db: db,
	}
}

// Create inserts a new action into the database
func (r *SQLiteActionRepository) Create(ctx context.Context, action *domain.Action) error {
	// Generate a new UUID if ID is not provided
	if action.ID == "" {
		action.ID = uuid.New().String()
	}

	now := time.Now()
	action.CreatedAt = now
	action.UpdatedAt = now

	query := `
		INSERT INTO actions (id, resource_id, name, description, attributes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.DB.ExecContext(
		ctx,
		query,
		action.ID,
		action.ResourceID,
		action.Name,
		action.Description,
		action.Attributes,
		action.CreatedAt,
		action.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create action: %w", err)
	}

	return nil
}

// GetByID retrieves an action by ID
func (r *SQLiteActionRepository) GetByID(ctx context.Context, id string) (*domain.Action, error) {
	query := `
		SELECT id, resource_id, name, description, attributes, created_at, updated_at
		FROM actions
		WHERE id = ?
	`

	var action domain.Action
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&action.ID,
		&action.ResourceID,
		&action.Name,
		&action.Description,
		&action.Attributes,
		&action.CreatedAt,
		&action.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("action not found")
		}
		return nil, fmt.Errorf("failed to get action: %w", err)
	}

	return &action, nil
}

// GetByResourceID retrieves actions by resource ID
func (r *SQLiteActionRepository) GetByResourceID(ctx context.Context, resourceID string) ([]*domain.Action, error) {
	query := `
		SELECT id, resource_id, name, description, attributes, created_at, updated_at
		FROM actions
		WHERE resource_id = ?
		ORDER BY created_at DESC
	`

	rows, err := r.db.DB.QueryContext(ctx, query, resourceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get actions: %w", err)
	}
	defer rows.Close()

	var actions []*domain.Action
	for rows.Next() {
		var action domain.Action
		err := rows.Scan(
			&action.ID,
			&action.ResourceID,
			&action.Name,
			&action.Description,
			&action.Attributes,
			&action.CreatedAt,
			&action.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan action: %w", err)
		}
		actions = append(actions, &action)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}

	return actions, nil
}

// List retrieves a paginated list of actions
func (r *SQLiteActionRepository) List(ctx context.Context, limit, offset int) ([]*domain.Action, error) {
	query := `
		SELECT id, resource_id, name, description, attributes, created_at, updated_at
		FROM actions
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list actions: %w", err)
	}
	defer rows.Close()

	var actions []*domain.Action
	for rows.Next() {
		var action domain.Action
		err := rows.Scan(
			&action.ID,
			&action.ResourceID,
			&action.Name,
			&action.Description,
			&action.Attributes,
			&action.CreatedAt,
			&action.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan action: %w", err)
		}
		actions = append(actions, &action)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}

	return actions, nil
}

// Update updates an existing action
func (r *SQLiteActionRepository) Update(ctx context.Context, action *domain.Action) error {
	action.UpdatedAt = time.Now()

	query := `
		UPDATE actions
		SET resource_id = ?, name = ?, description = ?, attributes = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.DB.ExecContext(
		ctx,
		query,
		action.ResourceID,
		action.Name,
		action.Description,
		action.Attributes,
		action.UpdatedAt,
		action.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update action: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("action not found")
	}

	return nil
}

// Delete deletes an action by ID
func (r *SQLiteActionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM actions WHERE id = ?`

	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete action: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("action not found")
	}

	return nil
}
