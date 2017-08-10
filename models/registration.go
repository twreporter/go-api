package models

import (
	"time"
)

// Registration this is bookmakrs table description
type Registration struct {
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
	Service       Service
	ServiceID     uint   `gorm:"primary_key" json:"service_id"`
	Email         string `gorm:"primary_key;size:100" json:"email"`
	User          User
	UserID        uint   `gorm:"default:0" json:"user_id"`
	Active        bool   `gorm:"default:0" json:"active"`
	ActivateToken string `gorm:"size:20" json:"active_token"`
}

// RegistrationJSON this is POST data in json format
type RegistrationJSON struct {
	Email         string `json:"email" binding:"required"`
	UserID        string `json:"uid"`
	Active        bool   `json:"active"`
	ActivateToken string `json:"active_token"`
}
