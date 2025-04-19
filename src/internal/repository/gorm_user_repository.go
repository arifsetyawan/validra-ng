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
type GormUserRepository struct {
	db *database.PostgresDB
}

// NewGormUserRepository creates a new GORM repository for users
func NewGormUserRepository(db *database.PostgresDB) domain.UserRepository {
	return &GormUserRepository{
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
}

// toDomain converts a GORM model to a domain model
func (u *User) toDomain() *domain.User {
	return &domain.User{
		ID:         u.ID,
		Username:   u.Username,
		Attributes: u.Attributes,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
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
	}
}

// Create inserts a new user into the database
func (r *GormUserRepository) Create(ctx context.Context, user *domain.User) error {
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
func (r *GormUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user User
	result := r.db.DB.WithContext(ctx).First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return user.toDomain(), nil
}

// GetByUsername retrieves a user by username
func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user User
	result := r.db.DB.WithContext(ctx).First(&user, "username = ?", username)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return user.toDomain(), nil
}

// List retrieves a paginated list of users
func (r *GormUserRepository) List(ctx context.Context, limit, offset int) ([]*domain.User, error) {
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
func (r *GormUserRepository) Update(ctx context.Context, user *domain.User) error {
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

// Delete removes a user from the database
func (r *GormUserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.DB.WithContext(ctx).Delete(&User{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
