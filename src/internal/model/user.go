package model

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	Username   string     `json:"username" gorm:"not null;unique"`
	Attributes []byte     `json:"attributes"` // JSON serialized attributes
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}
