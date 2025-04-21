package model

import (
	"time"
)

// UserSet represents a group of users defined by conditions
type UserSet struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Conditions  []byte     `json:"conditions"` // JSON serialized conditions
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}
