package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arifsetyawan/validra/src/internal/domain"
	"github.com/arifsetyawan/validra/src/pkg/database"
	"github.com/google/uuid"
)

// GormRoleRepository implements domain.RoleRepository using GORM with PostgreSQL
type RoleRepository struct {
	db *database.PostgresDB
}

// NewGormRoleRepository creates a new GORM repository for roles
func NewRoleRepository(db *database.PostgresDB) domain.RoleRepository {
	return &RoleRepository{
		db: db,
	}
}

// Role is the GORM model for roles
type Role struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
}

// toDomain converts a GORM model to a domain model
func (r *Role) toDomain() *domain.Role {
	return &domain.Role{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}
}

// fromDomain converts a domain model to a GORM model
func roleFromDomain(r *domain.Role) *Role {
	return &Role{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}
}

// Create inserts a new role into the database
func (r *RoleRepository) Create(ctx context.Context, role *domain.Role) error {
	// Generate a new UUID if not provided
	if role.ID == "" {
		role.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	gormRole := roleFromDomain(role)
	result := r.db.DB.WithContext(ctx).Create(gormRole)
	if result.Error != nil {
		return fmt.Errorf("failed to create role: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a role by ID
func (r *RoleRepository) GetByID(ctx context.Context, id string) (*domain.Role, error) {
	var role Role
	result := r.db.DB.WithContext(ctx).First(&role, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get role: %w", result.Error)
	}

	return role.toDomain(), nil
}

// List retrieves a paginated list of roles
func (r *RoleRepository) List(ctx context.Context, limit, offset int) ([]*domain.Role, error) {
	var roles []Role
	result := r.db.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&roles)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list roles: %w", result.Error)
	}

	domainRoles := make([]*domain.Role, len(roles))
	for i, role := range roles {
		domainRoles[i] = role.toDomain()
	}

	return domainRoles, nil
}

// Update updates a role in the database
func (r *RoleRepository) Update(ctx context.Context, role *domain.Role) error {
	// Update the timestamp
	role.UpdatedAt = time.Now()

	gormRole := roleFromDomain(role)
	result := r.db.DB.WithContext(ctx).Save(gormRole)
	if result.Error != nil {
		return fmt.Errorf("failed to update role: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

// Delete performs a soft delete on a role and returns the deleted role
func (r *RoleRepository) Delete(ctx context.Context, id string) (*domain.Role, error) {
	// First retrieve the role to return it after deletion
	var role Role
	getResult := r.db.DB.WithContext(ctx).First(&role, "id = ?", id)
	if getResult.Error != nil {
		return nil, fmt.Errorf("failed to get role: %w", getResult.Error)
	}

	// Perform soft delete
	now := time.Now()
	result := r.db.DB.WithContext(ctx).Model(&Role{}).Where("id = ?", id).Update("deleted_at", now)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to soft delete role: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("role not found")
	}

	// Update the retrieved role with deletion time
	role.DeletedAt = &now

	return role.toDomain(), nil
}
