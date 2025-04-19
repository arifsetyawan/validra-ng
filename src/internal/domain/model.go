package domain

import (
	"time"
)

// Resource represents an entity that can be accessed
type Resource struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Attributes  []byte    `json:"attributes"` // JSON serialized attributes
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Role represents a role in the system
type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// User represents a user in the system
type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username"`
	Attributes []byte    `json:"attributes"` // JSON serialized attributes
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// UserSet represents a group of users defined by conditions
type UserSet struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Conditions  []byte    `json:"conditions"` // JSON serialized conditions
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ResourceSet represents a group of resources defined by conditions
type ResourceSet struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Conditions  []byte    `json:"conditions"` // JSON serialized conditions
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Permission defines relationship between roles, users/user sets and resources/resource sets
type Permission struct {
	ID            string    `json:"id"`
	RoleID        string    `json:"role_id"`
	UserID        *string   `json:"user_id,omitempty"`
	UserSetID     *string   `json:"user_set_id,omitempty"`
	ResourceID    *string   `json:"resource_id,omitempty"`
	ResourceSetID *string   `json:"resource_set_id,omitempty"`
	Effect        string    `json:"effect"`               // "allow" or "deny"
	Conditions    []byte    `json:"conditions,omitempty"` // Additional conditions
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
