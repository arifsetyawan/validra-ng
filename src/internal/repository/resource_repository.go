package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/arifsetyawan/validra/src/internal/model"
	"github.com/arifsetyawan/validra/src/pkg/database"
	"github.com/google/uuid"
)

// GormResourceRepository implements domain.ResourceRepository using GORM with PostgreSQL
type ResourceRepository struct {
	db *database.PostgresDB
}

// ResourceRepository defines the methods for Resource data access
type ResourceRepositoryInterface interface {
	Create(ctx context.Context, resource *model.Resource) error
	GetByID(ctx context.Context, id string) (*model.Resource, error)
	List(ctx context.Context, limit, offset int) (*[]model.Resource, error)
	Update(ctx context.Context, resource *model.Resource) error
	Delete(ctx context.Context, id string) (*model.Resource, error)
}

// NewGormResourceRepository creates a new GORM repository for resources
func NewResourceRepository(db *database.PostgresDB) ResourceRepositoryInterface {
	return &ResourceRepository{
		db: db,
	}
}

// Create inserts a new resource into the database
func (r *ResourceRepository) Create(ctx context.Context, resource *model.Resource) error {
	// Generate a new UUID if ID is not provided
	if resource.ID == "" {
		resource.ID = uuid.New().String()
	}

	now := time.Now()
	resource.CreatedAt = now
	resource.UpdatedAt = now

	resourceModelObj := r.db.DB.WithContext(ctx).Model(&model.Resource{}).Preload("ResourceActions").Preload("ResourceABACOptions").Preload("ResourceReBACOptions")

	result := resourceModelObj.Create(resource)
	if result.Error != nil {
		return fmt.Errorf("failed to create resource: %w", result.Error)
	}

	return nil
}

// GetByID retrieves a resource by ID
func (r *ResourceRepository) GetByID(ctx context.Context, id string) (*model.Resource, error) {
	var resource model.Resource

	resourceModelObj := r.db.DB.WithContext(ctx).Model(&model.Resource{}).Preload("ResourceActions").Preload("ResourceABACOptions").Preload("ResourceReBACOptions")

	result := resourceModelObj.First(&resource, "id = ?", id)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get resource: %w", result.Error)
	}

	fmt.Println("Resource retrieved:", resource)

	return &resource, nil
}

// List retrieves a paginated list of resources
func (r *ResourceRepository) List(ctx context.Context, limit, offset int) (*[]model.Resource, error) {
	var resources []model.Resource

	resourceModelObj := r.db.DB.WithContext(ctx).Model(&model.Resource{}).Preload("ResourceActions").Preload("ResourceABACOptions").Preload("ResourceReBACOptions")

	result := resourceModelObj.Limit(limit).Offset(offset).Find(&resources)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to list resources: %w", result.Error)
	}

	return &resources, nil
}

// Update updates a resource in the database
func (r *ResourceRepository) Update(ctx context.Context, resource *model.Resource) error {
	resource.UpdatedAt = time.Now()

	result := r.db.DB.WithContext(ctx).Save(resource)
	if result.Error != nil {
		return fmt.Errorf("failed to update resource: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("resource not found")
	}

	return nil
}

// Delete performs a soft delete on a resource and returns the deleted resource
func (r *ResourceRepository) Delete(ctx context.Context, id string) (*model.Resource, error) {
	// First retrieve the resource to return it after deletion
	var resource model.Resource
	getResult := r.db.DB.WithContext(ctx).First(&resource, "id = ?", id)
	if getResult.Error != nil {
		return nil, fmt.Errorf("failed to get resource: %w", getResult.Error)
	}

	// Perform soft delete
	now := time.Now()
	result := r.db.DB.WithContext(ctx).Model(&resource).Where("id = ?", id).Update("deleted_at", now)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to soft delete resource: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("resource not found")
	}

	// Update the retrieved resource with deletion time
	resource.DeletedAt = &now

	return &resource, nil
}
