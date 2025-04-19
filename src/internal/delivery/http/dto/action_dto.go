package dto

import (
	"encoding/json"
	"time"

	"github.com/arifsetyawan/validra/src/internal/domain"
)

// CreateActionRequest represents the request payload for creating an action
type CreateActionRequest struct {
	ResourceID  string      `json:"resource_id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string      `json:"name" validate:"required" example:"read"`
	Description string      `json:"description" example:"Permission to read the resource"`
	Attributes  interface{} `json:"attributes,omitempty" swaggertype:"object"`
}

// UpdateActionRequest represents the request payload for updating an action
type UpdateActionRequest struct {
	ResourceID  string      `json:"resource_id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string      `json:"name" validate:"required" example:"read"`
	Description string      `json:"description" example:"Permission to read the resource"`
	Attributes  interface{} `json:"attributes,omitempty" swaggertype:"object"`
}

// ActionResponse represents the response model for an action
type ActionResponse struct {
	ID          string      `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ResourceID  string      `json:"resource_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string      `json:"name" example:"read"`
	Description string      `json:"description" example:"Permission to read the resource"`
	Attributes  interface{} `json:"attributes,omitempty" swaggertype:"object"`
	CreatedAt   time.Time   `json:"created_at" example:"2025-04-19T12:00:00Z"`
	UpdatedAt   time.Time   `json:"updated_at" example:"2025-04-19T12:00:00Z"`
}

// ListActionsResponse represents a paginated list of actions
type ListActionsResponse struct {
	Actions []ActionResponse `json:"actions"`
	Total   int              `json:"total" example:"10"`
}

// ToActionResponse converts a domain.Action to ActionResponse
func ToActionResponse(a *domain.Action) ActionResponse {
	var attributes interface{}
	// Only unmarshal if there are attributes
	if a.Attributes != nil && len(a.Attributes) > 0 {
		if err := json.Unmarshal(a.Attributes, &attributes); err != nil {
			// Fall back to raw bytes if unmarshaling fails
			attributes = a.Attributes
		}
	}

	return ActionResponse{
		ID:          a.ID,
		ResourceID:  a.ResourceID,
		Name:        a.Name,
		Description: a.Description,
		Attributes:  attributes,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

// ToActionDomain converts a CreateActionRequest to domain.Action
func (r *CreateActionRequest) ToActionDomain() *domain.Action {
	var attributesBytes []byte
	var err error

	// Only marshal if there are attributes
	if r.Attributes != nil {
		attributesBytes, err = json.Marshal(r.Attributes)
		if err != nil {
			// Log error or handle it accordingly
			attributesBytes = []byte("{}")
		}
	}

	return &domain.Action{
		ResourceID:  r.ResourceID,
		Name:        r.Name,
		Description: r.Description,
		Attributes:  attributesBytes,
	}
}

// UpdateActionDomain updates a domain.Action with values from UpdateActionRequest
func (r *UpdateActionRequest) UpdateActionDomain(action *domain.Action) {
	action.ResourceID = r.ResourceID
	action.Name = r.Name
	action.Description = r.Description

	// Only update attributes if provided
	if r.Attributes != nil {
		attributesBytes, err := json.Marshal(r.Attributes)
		if err == nil {
			action.Attributes = attributesBytes
		}
		// If error, we just don't update the attributes
	}
}
