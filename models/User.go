package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User ...
type User struct {
	gorm.Model              // contains fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	Email            string `gorm:"size:100"`
	FirstName        string `gorm:"size:50"`
	LastName         string `gorm:"size:50"`
	SecurityID       string `gorm:"size:20"`
	PassportID       string `gorm:"size:30"`
	City             string `gorm:"size:45"`
	State            string `gorm:"size:45"`
	Country          string `gorm:"size:45"`
	Zip              string `gorm:"size:20"`
	Address          string
	Phone            string `gorm:"size:20"`
	Privilege        int    `gorm:"size:2"`
	RegistrationDate time.Time
	Birthday         time.Time
	Gender           string `gorm:"size:2"`  // e.g., "M", "F" ...
	Education        string `gorm:"size:20"` // e.g., "High School"
	EnableEmail      int    `gorm:"size:2"`
}
