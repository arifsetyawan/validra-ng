package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arifsetyawan/validra/src/internal/domain"
	"github.com/arifsetyawan/validra/src/pkg/database"
	"github.com/google/uuid"
)

// GormUserRepository implements domain.UserRepository using GORM with PostgreSQL
type UserRepository struct {
	db *database.PostgresDB
}

// NewUserRepository creates a new GORM repository for users
func NewUserRepository(db *database.PostgresDB) domain.UserRepository {
	return &UserRepository{
		db: db,
	}
}

// User is the GORM model for users
type User struct {
	ID         string `gorm:"primaryKey"`
	Username   string `gorm:"not null;unique"`
	Attributes []byte
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`
}

// toDomain converts a GORM model to a domain model
func (u *User) toDomain() *domain.User {
	return &domain.User{
		ID:         u.ID,
		Username:   u.Username,
		Attributes: u.Attributes,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
		DeletedAt:  u.DeletedAt,
	}
}

// fromDomain converts a domain model to a GORM model
func userFromDomain(u *domain.User) *User {
	return &User{
		ID:         u.ID,
		Username:   u.Username,
		Attributes: u.Attributes,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
		DeletedAt:  u.DeletedAt,
	}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	// Generate a new UUID if ID is not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	gormUser := userFromDomain(user)
	result := r.db.DB.WithContext(ctx).Create(gormUser)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user User
	result := r.db.DB.WithContext(ctx).First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return user.toDomain(), nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user User
	result := r.db.DB.WithContext(ctx).First(&user, "username = ?", username)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return user.toDomain(), nil
}

// List retrieves a paginated list of users
func (r *UserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	var users []User
	result := r.db.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list users: %w", result.Error)
	}

	domainUsers := make([]*domain.User, len(users))
	for i, user := range users {
		domainUsers[i] = user.toDomain()
	}

	return domainUsers, nil
}

// Update updates a user in the database
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	user.UpdatedAt = time.Now()

	gormUser := userFromDomain(user)
	result := r.db.DB.WithContext(ctx).Save(gormUser)
	if result.Error != nil {
		return fmt.Errorf("failed to update user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Delete performs a soft delete on a user and returns the deleted user
func (r *UserRepository) Delete(ctx context.Context, id string) (*domain.User, error) {
	// First retrieve the user to return it after deletion
	var user User
	getResult := r.db.DB.WithContext(ctx).First(&user, "id = ?", id)
	if getResult.Error != nil {
		return nil, fmt.Errorf("failed to get user: %w", getResult.Error)
	}

	// Perform soft delete
	now := time.Now()
	result := r.db.DB.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("deleted_at", now)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to soft delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found")
	}

	// Update the retrieved user with deletion time
	user.DeletedAt = &now

	return user.toDomain(), nil
}
