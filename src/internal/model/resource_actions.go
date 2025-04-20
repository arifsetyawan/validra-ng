package model

import (
	"time"
)

// ResourceActions is the GORM model for actions
type ResourceActions struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	ResourceID string     `json:"resource_id" gorm:"index"`
	Name       string     `json:"name" gorm:"not null"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"index"`
}
