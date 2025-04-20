package model

import (
	"time"
)

type ResourceABACOptions struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	ResourceID string     `json:"resource_id" gorm:"index"`
	Name       string     `json:"name" gorm:"not null"`
	Type       string     `json:"type" gorm:"not null"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"index"`
}
