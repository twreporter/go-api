package models

import (
	"github.com/jinzhu/gorm"
)

// Service this is service table description
type Service struct {
	gorm.Model
	Name string `gorm:"size:100;unique_index;not null"`
}

type ServiceJSON struct {
	ID   uint   `json:"id"`
	Name string `json:"name" binding:"required"`
}
