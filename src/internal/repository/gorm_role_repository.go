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
type GormRoleRepository struct {
	db *database.PostgresDB
}

// NewGormRoleRepository creates a new GORM repository for roles
func NewGormRoleRepository(db *database.PostgresDB) domain.RoleRepository {
	return &GormRoleRepository{
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
}

// toDomain converts a GORM model to a domain model
func (r *Role) toDomain() *domain.Role {
	return &domain.Role{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
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
	}
}

// Create inserts a new role into the database
func (r *GormRoleRepository) Create(ctx context.Context, role *domain.Role) error {
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
func (r *GormRoleRepository) GetByID(ctx context.Context, id string) (*domain.Role, error) {
	var role Role
	result := r.db.DB.WithContext(ctx).First(&role, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get role: %w", result.Error)
	}

	return role.toDomain(), nil
}

// List retrieves a paginated list of roles
func (r *GormRoleRepository) List(ctx context.Context, limit, offset int) ([]*domain.Role, error) {
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
func (r *GormRoleRepository) Update(ctx context.Context, role *domain.Role) error {
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

// Delete removes a role from the database
func (r *GormRoleRepository) Delete(ctx context.Context, id string) error {
	result := r.db.DB.WithContext(ctx).Delete(&Role{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete role: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}
