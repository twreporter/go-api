package models

import (
	"time"
)

// Registration this is bookmakrs table description
type Registration struct {
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
	Service       Service
	ServiceID     uint   `gorm:"primary_key"`
	Email         string `gorm:"primary_key;size:100"`
	UID           uint
	Active        bool   `gorm:"default:0"`
	ActivateToken string `gorm:"size:20"`
}

// RegistrationJSON this is POST data in json format
type RegistrationJSON struct {
	Email         string `json:"email" binding:"required"`
	Service       string `json:"service"`
	UID           uint   `json:"uid"`
	Active        bool   `json:"active"`
	ActivateToken string `json:"active_token"`
}
