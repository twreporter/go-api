package models

import (
	"time"
)

// Service this is service table description
type Service struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Name      string     `gorm:"size:100;unique_index;not null" json:"name"`
}

// ServiceJSON ...
type ServiceJSON struct {
	ID   uint   `json:"id"`
	Name string `json:"name" binding:"required"`
}
