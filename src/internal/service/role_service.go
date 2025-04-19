package service

import (
	"context"
	"fmt"

	"github.com/arifsetyawan/validra/src/internal/domain"
)

// RoleService handles business logic for roles
type RoleService struct {
	roleRepo domain.RoleRepository
}

// NewRoleService creates a new RoleService
func NewRoleService(roleRepo domain.RoleRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
	}
}

// CreateRole creates a new role
func (s *RoleService) CreateRole(ctx context.Context, role *domain.Role) error {
	if role.Name == "" {
		return fmt.Errorf("role name is required")
	}

	return s.roleRepo.Create(ctx, role)
}

// GetRoleByID retrieves a role by ID
func (s *RoleService) GetRoleByID(ctx context.Context, id string) (*domain.Role, error) {
	return s.roleRepo.GetByID(ctx, id)
}

// ListRoles retrieves a paginated list of roles
func (s *RoleService) ListRoles(ctx context.Context, limit, offset int) ([]*domain.Role, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	return s.roleRepo.List(ctx, limit, offset)
}

// UpdateRole updates an existing role
func (s *RoleService) UpdateRole(ctx context.Context, role *domain.Role) error {
	if role.Name == "" {
		return fmt.Errorf("role name is required")
	}

	// Check if the role exists
	_, err := s.roleRepo.GetByID(ctx, role.ID)
	if err != nil {
		return fmt.Errorf("role not found")
	}

	return s.roleRepo.Update(ctx, role)
}

// DeleteRole deletes a role by ID
func (s *RoleService) DeleteRole(ctx context.Context, id string) error {
	// Check if the role exists
	_, err := s.roleRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("role not found")
	}

	return s.roleRepo.Delete(ctx, id)
}
