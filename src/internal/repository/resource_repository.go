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
type ResourceRepository struct {
	db *database.PostgresDB
}

// ResourceRepository defines the methods for Resource data access
type ResourceRepositoryInterface interface {
	Create(ctx context.Context, resource *domain.Resource) error
	GetByID(ctx context.Context, id string) (*domain.Resource, error)
	List(ctx context.Context, limit, offset int) ([]*domain.Resource, error)
	Update(ctx context.Context, resource *domain.Resource) error
	Delete(ctx context.Context, id string) (*domain.Resource, error)
}

// NewGormResourceRepository creates a new GORM repository for resources
func NewResourceRepository(db *database.PostgresDB) ResourceRepositoryInterface {
	return &ResourceRepository{
		db: db,
	}
}

// Resource is the GORM model for resources
type Resource struct {
	ID          string `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string
	Attributes  []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `gorm:"index"`
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
		DeletedAt:   r.DeletedAt,
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
		DeletedAt:   r.DeletedAt,
	}
}

// Create inserts a new resource into the database
func (r *ResourceRepository) Create(ctx context.Context, resource *domain.Resource) error {
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
func (r *ResourceRepository) GetByID(ctx context.Context, id string) (*domain.Resource, error) {
	var resource Resource
	result := r.db.DB.WithContext(ctx).First(&resource, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get resource: %w", result.Error)
	}

	return resource.toDomain(), nil
}

// List retrieves a paginated list of resources
func (r *ResourceRepository) List(ctx context.Context, limit, offset int) ([]*domain.Resource, error) {
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
func (r *ResourceRepository) Update(ctx context.Context, resource *domain.Resource) error {
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

// Delete performs a soft delete on a resource and returns the deleted resource
func (r *ResourceRepository) Delete(ctx context.Context, id string) (*domain.Resource, error) {
	// First retrieve the resource to return it after deletion
	var resource Resource
	getResult := r.db.DB.WithContext(ctx).First(&resource, "id = ?", id)
	if getResult.Error != nil {
		return nil, fmt.Errorf("failed to get resource: %w", getResult.Error)
	}

	// Perform soft delete
	now := time.Now()
	result := r.db.DB.WithContext(ctx).Model(&Resource{}).Where("id = ?", id).Update("deleted_at", now)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to soft delete resource: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("resource not found")
	}

	// Update the retrieved resource with deletion time
	resource.DeletedAt = &now

	return resource.toDomain(), nil
}
