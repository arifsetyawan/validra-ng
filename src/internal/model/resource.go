package model

import (
	"time"
)

// Resource is the GORM model for resources
type Resource struct {
	ID                   string                 `json:"id" gorm:"primaryKey"`
	Name                 string                 `json:"name" gorm:"not null"`
	Description          string                 `json:"description"`
	ResourceABACOptions  []ResourceABACOptions  `json:"abac_options"`
	ResourceReBACOptions []ResourceReBACOptions `json:"rebac_options"`
	CreatedAt            time.Time              `json:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at"`
	DeletedAt            *time.Time             `json:"deleted_at" gorm:"index"`
}

type ResourceABACOptions struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	ResourceID string     `json:"resource_id" gorm:"index"`
	Name       string     `json:"name" gorm:"not null"`
	Type       string     `json:"type" gorm:"not null"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"index"`
}

type ResourceReBACOptions struct {
	ID           string     `json:"id" gorm:"primaryKey"`
	ResourceID   string     `json:"resource_id" gorm:"index"`
	Name         string     `json:"name" gorm:"not null"`
	Relationship string     `json:"relationship" gorm:"not null"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`
}
