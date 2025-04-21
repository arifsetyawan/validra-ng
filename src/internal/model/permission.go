package model

import "time"

type Permission struct {
	ID            string     `json:"id"`
	RoleID        string     `json:"role_id"`
	UserID        *string    `json:"user_id,omitempty"`
	UserSetID     *string    `json:"user_set_id,omitempty"`
	ResourceID    *string    `json:"resource_id,omitempty"`
	ResourceSetID *string    `json:"resource_set_id,omitempty"`
	Effect        string     `json:"effect"`               // "allow" or "deny"
	Conditions    []byte     `json:"conditions,omitempty"` // Additional conditions
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}
