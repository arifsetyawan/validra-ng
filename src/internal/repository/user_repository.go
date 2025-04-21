package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arifsetyawan/validra/src/internal/model"
	"github.com/arifsetyawan/validra/src/pkg/database"
	"github.com/google/uuid"
)

// GormUserRepository implements domain.UserRepository using GORM with PostgreSQL
type UserRepository struct {
	db *database.PostgresDB
}

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	List(ctx context.Context, limit, offset int) (*[]model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) (*model.User, error)
}

// NewUserRepository creates a new GORM repository for users
func NewUserRepository(db *database.PostgresDB) UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	// Generate a new UUID if ID is not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result := r.db.DB.WithContext(ctx).Create(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	result := r.db.DB.WithContext(ctx).First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	result := r.db.DB.WithContext(ctx).First(&user, "username = ?", username)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

// List retrieves a paginated list of users
func (r *UserRepository) List(ctx context.Context, limit, offset int) (*[]model.User, error) {
	var users []model.User
	result := r.db.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list users: %w", result.Error)
	}

	return &users, nil
}

// Update updates a user in the database
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	user.UpdatedAt = time.Now()

	result := r.db.DB.WithContext(ctx).Save(user)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete performs a soft delete on a user and returns the deleted user
func (r *UserRepository) Delete(ctx context.Context, id string) (*model.User, error) {
	// First retrieve the user to return it after deletion
	var user model.User
	getResult := r.db.DB.WithContext(ctx).First(&user, "id = ?", id)
	if getResult.Error != nil {
		return nil, fmt.Errorf("failed to get user: %w", getResult.Error)
	}

	// Perform soft delete
	now := time.Now()
	result := r.db.DB.WithContext(ctx).Model(&user).Where("id = ?", id).Update("deleted_at", now)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to soft delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}

	// Update the retrieved user with deletion time
	user.DeletedAt = &now

	return &user, nil
}
