package dto

import (
	"encoding/json"
	"time"

	"github.com/arifsetyawan/validra/src/internal/domain"
)

// CreateResourceRequest represents the request payload for creating a resource
type CreateResourceRequest struct {
	Name        string      `json:"name" validate:"required" example:"Sample Resource"`
	Description string      `json:"description" example:"This is a sample resource description"`
	Attributes  interface{} `json:"attributes,omitempty" swaggertype:"object"`
}

// UpdateResourceRequest represents the request payload for updating a resource
type UpdateResourceRequest struct {
	Name        string      `json:"name" validate:"required" example:"Updated Resource"`
	Description string      `json:"description" example:"This is an updated resource description"`
	Attributes  interface{} `json:"attributes,omitempty" swaggertype:"object"`
}

// ResourceResponse represents the response model for a resource
type ResourceResponse struct {
	ID          string      `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name        string      `json:"name" example:"Sample Resource"`
	Description string      `json:"description" example:"This is a sample resource description"`
	Attributes  interface{} `json:"attributes,omitempty" swaggertype:"object"`
	CreatedAt   time.Time   `json:"created_at" example:"2025-04-19T12:00:00Z"`
	UpdatedAt   time.Time   `json:"updated_at" example:"2025-04-19T12:00:00Z"`
}

// ListResourcesResponse represents a paginated list of resources
type ListResourcesResponse struct {
	Resources []ResourceResponse `json:"resources"`
	Total     int                `json:"total" example:"10"`
}

// ToResourceResponse converts a domain.Resource to ResourceResponse
func ToResourceResponse(r *domain.Resource) ResourceResponse {
	var attributes interface{}
	// Only unmarshal if there are attributes
	if r.Attributes != nil && len(r.Attributes) > 0 {
		if err := json.Unmarshal(r.Attributes, &attributes); err != nil {
			// Fall back to raw bytes if unmarshaling fails
			attributes = r.Attributes
		}
	}

	return ResourceResponse{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Attributes:  attributes,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// ToResourceDomain converts a CreateResourceRequest to domain.Resource
func (r *CreateResourceRequest) ToResourceDomain() *domain.Resource {
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

	return &domain.Resource{
		Name:        r.Name,
		Description: r.Description,
		Attributes:  attributesBytes,
	}
}

// UpdateResourceDomain updates a domain.Resource with values from UpdateResourceRequest
func (r *UpdateResourceRequest) UpdateResourceDomain(resource *domain.Resource) {
	resource.Name = r.Name
	resource.Description = r.Description

	// Only update attributes if provided
	if r.Attributes != nil {
		attributesBytes, err := json.Marshal(r.Attributes)
		if err == nil {
			resource.Attributes = attributesBytes
		}
		// If error, we just don't update the attributes
	}
}
