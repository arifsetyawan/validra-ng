package model

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	Email      string     `json:"email" gorm:"not null;unique"`
	Entity     string     `json:"entity" gorm:"not null"`
	Attributes []byte     `json:"attributes"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}
