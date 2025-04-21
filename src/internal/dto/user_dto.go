package dto

import (
	"encoding/json"
	"time"

	"github.com/arifsetyawan/validra/src/internal/model"
)

// CreateUserRequest represents the request payload for creating a user
type CreateUserRequest struct {
	Username   string      `json:"username" validate:"required" example:"john_doe"`
	Attributes interface{} `json:"attributes,omitempty" swaggertype:"object"`
}

// UpdateUserRequest represents the request payload for updating a user
type UpdateUserRequest struct {
	Username   string      `json:"username" validate:"required" example:"john_smith"`
	Attributes interface{} `json:"attributes,omitempty" swaggertype:"object"`
}

// UserResponse represents the response model for a user
type UserResponse struct {
	ID         string      `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Username   string      `json:"username" example:"john_doe"`
	Attributes interface{} `json:"attributes,omitempty" swaggertype:"object"`
	CreatedAt  time.Time   `json:"created_at" example:"2025-04-19T12:00:00Z"`
	UpdatedAt  time.Time   `json:"updated_at" example:"2025-04-19T12:00:00Z"`
}

// ListUsersResponse represents a paginated list of users
type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int            `json:"total" example:"10"`
}

// ToUserResponse converts a model.User to UserResponse
func ToUserResponse(u *model.User) UserResponse {
	var attributes interface{}
	if u.Attributes != nil && len(u.Attributes) > 0 {
		json.Unmarshal(u.Attributes, &attributes)
	}

	return UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		Attributes: attributes,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

// ToUserDomain converts a CreateUserRequest to domain.User
func (r *CreateUserRequest) ToUserDomain() *model.User {
	var attributesBytes []byte
	if r.Attributes != nil {
		attributesBytes, _ = json.Marshal(r.Attributes)
	}

	return &model.User{
		Username:   r.Username,
		Attributes: attributesBytes,
	}
}

// UpdateUserDomain updates a domain.User with values from UpdateUserRequest
func (r *UpdateUserRequest) UpdateUserDomain(user *model.User) {
	user.Username = r.Username

	if r.Attributes != nil {
		attributesBytes, _ := json.Marshal(r.Attributes)
		user.Attributes = attributesBytes
	}
}
