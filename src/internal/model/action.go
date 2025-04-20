package model

import (
	"time"
)

// Action is the GORM model for actions
type Action struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	ResourceID  string     `json:"resource_id" gorm:"not null;index"`
	Resource    Resource   `json:"resource" gorm:"foreignKey:ResourceID"`
	Name        string     `json:"name" gorm:"not null"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}
