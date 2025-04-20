package dto

import (
	"time"

	"github.com/arifsetyawan/validra/src/internal/model"
)

// CreateRoleRequest is the DTO for creating a new role
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// UpdateRoleRequest is the DTO for updating an existing role
type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// RoleResponse is the DTO for role responses
type RoleResponse struct {
	ID          string    `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string    `json:"name" example:"admin"`
	Description string    `json:"description" example:"Administrator role with full access"`
	CreatedAt   time.Time `json:"created_at" example:"2025-04-19T12:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2025-04-19T12:00:00Z"`
}

// ListRolesResponse is the DTO for listing roles
type ListRolesResponse struct {
	Roles []RoleResponse `json:"roles"`
	Total int            `json:"total" example:"10"`
}

// ToRoleDomain converts a CreateRoleRequest to model.Role
func (r *CreateRoleRequest) ToRoleDomain() *model.Role {
	return &model.Role{
		Name:        r.Name,
		Description: r.Description,
	}
}

// UpdateRoleDomain updates a model.Role with values from UpdateRoleRequest
func (r *UpdateRoleRequest) UpdateRoleDomain(role *model.Role) {
	role.Name = r.Name
	role.Description = r.Description
}

// ToRoleResponse converts a domain.Role to RoleResponse
func ToRoleResponse(role *model.Role) RoleResponse {
	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}
