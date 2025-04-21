package dto

import (
	"encoding/json"
	"time"

	"github.com/arifsetyawan/validra/src/internal/model"
)

// CreateResourceRequest represents the request payload for creating a resource
type CreateResourceRequest struct {
	Name                 string                   `json:"name" validate:"required" example:"Sample Resource"`
	Description          string                   `json:"description" example:"This is a sample resource description"`
	Actions              []string                 `json:"actions" example:"[\"read\", \"write\"]"`
	ResourceABACOptions  []map[string]interface{} `json:"abac_options,omitempty" swaggertype:"array,object"`
	ResourceReBACOptions []map[string]interface{} `json:"rebac_options,omitempty" swaggertype:"array,object"`
}

// UpdateResourceRequest represents the request payload for updating a resource
type UpdateResourceRequest struct {
	ID                   string        `json:"id" validate:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name                 string        `json:"name" validate:"required" example:"Updated Resource"`
	Description          string        `json:"description" example:"This is an updated resource description"`
	Actions              []interface{} `json:"actions" swaggertype:"array,object"`
	ResourceABACOptions  []interface{} `json:"abac_options,omitempty" swaggertype:"array,object"`
	ResourceReBACOptions []interface{} `json:"rebac_options,omitempty" swaggertype:"array,object"`
}

// ResourceResponse represents the response model for a resource
type ResourceResponse struct {
	ID                   string        `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Name                 string        `json:"name" example:"Sample Resource"`
	Description          string        `json:"description" example:"This is a sample resource description"`
	Actions              []interface{} `json:"actions" swaggertype:"array,object"`
	ResourceABACOptions  []interface{} `json:"abac_options,omitempty" swaggertype:"array,object"`
	ResourceReBACOptions []interface{} `json:"rebac_options,omitempty" swaggertype:"array,object"`
	CreatedAt            time.Time     `json:"created_at" example:"2025-04-19T12:00:00Z"`
	UpdatedAt            time.Time     `json:"updated_at" example:"2025-04-19T12:00:00Z"`
}

// ListResourcesResponse represents a paginated list of resources
type ListResourcesResponse struct {
	Resources []ResourceResponse `json:"resources"`
	Total     int                `json:"total" example:"10"`
}

// ToResourceResponse converts a model.Resource to ResourceResponse
func ToResourceResponse(r *model.Resource) ResourceResponse {

	var actions []interface{}
	var abacOptions []interface{}
	var rebacOptions []interface{}

	// Convert to []interface{}
	for _, action := range r.ResourceActions {
		// Convert struct to map using JSON marshaling and unmarshaling
		actionBytes, _ := json.Marshal(action)
		var actionMap []interface{}
		_ = json.Unmarshal(actionBytes, &actionMap)
		actions = append(actions, actionMap)
	}

	for _, option := range r.ResourceABACOptions {
		// Convert struct to map using JSON marshaling and unmarshaling
		optionBytes, _ := json.Marshal(option)
		var optionMap []interface{}
		_ = json.Unmarshal(optionBytes, &optionMap)
		abacOptions = append(abacOptions, optionMap)
	}

	for _, option := range r.ResourceReBACOptions {
		// Convert struct to map using JSON marshaling and unmarshaling
		optionBytes, _ := json.Marshal(option)
		var optionMap []interface{}
		_ = json.Unmarshal(optionBytes, &optionMap)
		rebacOptions = append(rebacOptions, optionMap)
	}

	return ResourceResponse{
		ID:                   r.ID,
		Name:                 r.Name,
		Description:          r.Description,
		Actions:              actions,
		ResourceABACOptions:  abacOptions,
		ResourceReBACOptions: rebacOptions,
		CreatedAt:            r.CreatedAt,
		UpdatedAt:            r.UpdatedAt,
	}
}

// ToResourceDomain converts a CreateResourceRequest to model.Resource
func (r *CreateResourceRequest) ToResourceModel() *model.Resource {

	// Convert actions to []model.Action
	var actions []model.ResourceActions
	for _, action := range r.Actions {
		var actionModel model.ResourceActions
		jsonData, _ := json.Marshal(action)
		json.Unmarshal(jsonData, &actionModel)
		actions = append(actions, actionModel)
	}

	// Convert r.ResourceABACOptions to model.Resource.ResourceABACOptions
	var abacOptions []model.ResourceABACOptions
	for _, value := range r.ResourceABACOptions {
		var abacOption model.ResourceABACOptions
		jsonData, _ := json.Marshal(value)
		json.Unmarshal(jsonData, &abacOption)
		abacOptions = append(abacOptions, abacOption)
	}

	// Convert r.ResourceReBACOptions to model.Resource.ResourceReBACOptions
	var rebacOptions []model.ResourceReBACOptions
	for _, value := range r.ResourceReBACOptions {
		var rebacOption model.ResourceReBACOptions
		jsonData, _ := json.Marshal(value)
		json.Unmarshal(jsonData, &rebacOption)
		rebacOptions = append(rebacOptions, rebacOption)
	}

	return &model.Resource{
		Name:                 r.Name,
		Description:          r.Description,
		ResourceActions:      actions,
		ResourceABACOptions:  abacOptions,
		ResourceReBACOptions: rebacOptions,
	}
}
