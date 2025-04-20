package service

import (
	"context"
	"fmt"

	"github.com/arifsetyawan/validra/src/internal/model"
	"github.com/arifsetyawan/validra/src/internal/repository"
)

// ResourceService handles business logic for resources
type ResourceService struct {
	resourceRepo repository.ResourceRepositoryInterface
}

// NewResourceService creates a new ResourceService
func NewResourceService(resourceRepo repository.ResourceRepositoryInterface) *ResourceService {
	return &ResourceService{
		resourceRepo: resourceRepo,
	}
}

// CreateResource creates a new resource
func (s *ResourceService) CreateResource(ctx context.Context, resource *model.Resource) error {
	if resource.Name == "" {
		return fmt.Errorf("resource name is required")
	}

	return s.resourceRepo.Create(ctx, resource)
}

// GetResourceByID retrieves a resource by ID
func (s *ResourceService) GetResourceByID(ctx context.Context, id string) (*model.Resource, error) {
	return s.resourceRepo.GetByID(ctx, id)
}

// ListResources retrieves a paginated list of resources
func (s *ResourceService) ListResources(ctx context.Context, limit, offset int) (*[]model.Resource, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	return s.resourceRepo.List(ctx, limit, offset)
}

// UpdateResource updates an existing resource
func (s *ResourceService) UpdateResource(ctx context.Context, resource *model.Resource) error {
	if resource.ID == "" {
		return fmt.Errorf("resource ID is required")
	}
	if resource.Name == "" {
		return fmt.Errorf("resource name is required")
	}

	return s.resourceRepo.Update(ctx, resource)
}

// DeleteResource deletes a resource by ID
func (s *ResourceService) DeleteResource(ctx context.Context, id string) (*model.Resource, error) {
	return s.resourceRepo.Delete(ctx, id)
}

// ResourceRepository returns the resource repository
func (s *ResourceService) ResourceRepository() repository.ResourceRepositoryInterface {
	return s.resourceRepo
}
