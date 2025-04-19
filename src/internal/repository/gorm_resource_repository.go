package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arifsetyawan/validra/src/internal/domain"
	"github.com/arifsetyawan/validra/src/pkg/database"
	"github.com/google/uuid"
)

// GormResourceRepository implements domain.ResourceRepository using GORM with PostgreSQL
type GormResourceRepository struct {
	db *database.PostgresDB
}

// NewGormResourceRepository creates a new GORM repository for resources
func NewGormResourceRepository(db *database.PostgresDB) domain.ResourceRepository {
	return &GormResourceRepository{
		db: db,
	}
}

// Resource is the GORM model for resources
type Resource struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Attributes  []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// toDomain converts a GORM model to a domain model
func (r *Resource) toDomain() *domain.Resource {
	return &domain.Resource{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Attributes:  r.Attributes,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// fromDomain converts a domain model to a GORM model
func fromDomain(r *domain.Resource) *Resource {
	return &Resource{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Attributes:  r.Attributes,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// Create inserts a new resource into the database
func (r *GormResourceRepository) Create(ctx context.Context, resource *domain.Resource) error {
	// Generate a new UUID if ID is not provided
	if resource.ID == "" {
		resource.ID = uuid.New().String()
	}

	now := time.Now()
	resource.CreatedAt = now
	resource.UpdatedAt = now

	gormResource := fromDomain(resource)
	result := r.db.DB.WithContext(ctx).Create(gormResource)
	if result.Error != nil {
		return fmt.Errorf("failed to create resource: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a resource by ID
func (r *GormResourceRepository) GetByID(ctx context.Context, id string) (*domain.Resource, error) {
	var resource Resource
	result := r.db.DB.WithContext(ctx).First(&resource, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get resource: %w", result.Error)
	}

	return resource.toDomain(), nil
}

// List retrieves a paginated list of resources
func (r *GormResourceRepository) List(ctx context.Context, limit, offset int) ([]*domain.Resource, error) {
	var resources []Resource
	result := r.db.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&resources)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list resources: %w", result.Error)
	}

	domainResources := make([]*domain.Resource, len(resources))
	for i, resource := range resources {
		domainResources[i] = resource.toDomain()
	}

	return domainResources, nil
}

// Update updates a resource in the database
func (r *GormResourceRepository) Update(ctx context.Context, resource *domain.Resource) error {
	resource.UpdatedAt = time.Now()

	gormResource := fromDomain(resource)
	result := r.db.DB.WithContext(ctx).Save(gormResource)
	if result.Error != nil {
		return fmt.Errorf("failed to update resource: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("resource not found")
	}

	return nil
}

// Delete removes a resource from the database
func (r *GormResourceRepository) Delete(ctx context.Context, id string) error {
	result := r.db.DB.WithContext(ctx).Delete(&Resource{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete resource: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("resource not found")
	}

	return nil
}
