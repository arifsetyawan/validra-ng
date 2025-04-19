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

// SQLiteUserRepository implements domain.UserRepository using SQLite
type SQLiteUserRepository struct {
	db *database.SQLiteDB
}

// NewSQLiteUserRepository creates a new SQLite repository for users
func NewSQLiteUserRepository(db *database.SQLiteDB) domain.UserRepository {
	return &SQLiteUserRepository{
		db: db,
	}
}

// Create inserts a new user into the database
func (r *SQLiteUserRepository) Create(ctx context.Context, user *domain.User) error {
	// Generate a new UUID if ID is not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	query := `
		INSERT INTO users (id, username, attributes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.DB.ExecContext(
		ctx,
		query,
		user.ID,
		user.Username,
		user.Attributes,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *SQLiteUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, username, attributes, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user domain.User
	err := r.db.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Attributes,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *SQLiteUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, attributes, created_at, updated_at
		FROM users
		WHERE username = ?
	`

	var user domain.User
	err := r.db.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Attributes,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// List retrieves a paginated list of users
func (r *SQLiteUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	query := `
		SELECT id, username, attributes, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Attributes,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error in rows: %w", err)
	}

	return users, nil
}

// Update updates an existing user
func (r *SQLiteUserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET username = ?, attributes = ?, updated_at = ?
		WHERE id = ?
	`

	result, err := r.db.DB.ExecContext(
		ctx,
		query,
		user.Username,
		user.Attributes,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete deletes a user by ID
func (r *SQLiteUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
