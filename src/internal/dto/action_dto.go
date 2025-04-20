package dto

import (
	"time"

	"github.com/arifsetyawan/validra/src/internal/model"
)

// CreateActionRequest represents the request payload for creating an action
type CreateActionRequest struct {
	ResourceID  string `json:"resource_id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string `json:"name" validate:"required" example:"read"`
	Description string `json:"description" example:"Permission to read the resource"`
}

// UpdateActionRequest represents the request payload for updating an action
type UpdateActionRequest struct {
	ResourceID  string `json:"resource_id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string `json:"name" validate:"required" example:"read"`
	Description string `json:"description" example:"Permission to read the resource"`
}

// ActionResponse represents the response model for an action
type ActionResponse struct {
	ID          string           `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ResourceID  string           `json:"resource_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string           `json:"name" example:"read"`
	Description string           `json:"description" example:"Permission to read the resource"`
	CreatedAt   time.Time        `json:"created_at" example:"2025-04-19T12:00:00Z"`
	UpdatedAt   time.Time        `json:"updated_at" example:"2025-04-19T12:00:00Z"`
	Resource    ResourceResponse `json:"resource,omitempty"`
}

// ListActionsResponse represents a paginated list of actions
type ListActionsResponse struct {
	Actions []ActionResponse `json:"actions"`
	Total   int              `json:"total" example:"10"`
}

// ToActionResponse converts a domain.Action to ActionResponse
func ToActionResponse(a *model.Action) ActionResponse {
	return ActionResponse{
		ID:          a.ID,
		ResourceID:  a.ResourceID,
		Name:        a.Name,
		Description: a.Description,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// ToActionResponse converts a domain.Action to ActionResponse
func ToActionWithResourceResponse(a *model.Action) ActionResponse {

	// Convert resource to ResourceResponse
	resourceResponse := ResourceResponse{
		ID:          a.Resource.ID,
		Name:        a.Resource.Name,
		Description: a.Resource.Description,
		CreatedAt:   a.Resource.CreatedAt,
		UpdatedAt:   a.Resource.UpdatedAt,
	}

	return ActionResponse{
		ID:          a.ID,
		ResourceID:  a.ResourceID,
		Name:        a.Name,
		Description: a.Description,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		Resource:    resourceResponse,
	}
}

// ToActionDomain converts a CreateActionRequest to domain.Action
func (r *CreateActionRequest) ToActionDomain() *model.Action {
	return &model.Action{
		ResourceID:  r.ResourceID,
		Name:        r.Name,
		Description: r.Description,
	}
}

// UpdateActionDomain updates a domain.Action with values from UpdateActionRequest
func (r *UpdateActionRequest) UpdateActionDomain(action *model.Action) {
	action.ResourceID = r.ResourceID
	action.Name = r.Name
	action.Description = r.Description
}
