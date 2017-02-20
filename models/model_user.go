package models

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// User ...
type User struct {
	gorm.Model
	OAuthAccounts    []OAuthAccount  `gorm:"ForeignKey:UserID"` // a user has multiple oauth accounts //
	ReporterAccount  ReporterAccount `gorm:"ForeignKey:UserID"`
	Email            sql.NullString  `gorm:"size:100"`
	FirstName        sql.NullString  `gorm:"size:50"`
	LastName         sql.NullString  `gorm:"size:50"`
	SecurityID       sql.NullString  `gorm:"size:20"`
	PassportID       sql.NullString  `gorm:"size:30"`
	City             sql.NullString  `gorm:"size:45"`
	State            sql.NullString  `gorm:"size:45"`
	Country          sql.NullString  `gorm:"size:45"`
	Zip              sql.NullString  `gorm:"size:20"`
	Address          sql.NullString
	Phone            sql.NullString `gorm:"size:20"`
	Privilege        sql.NullInt64  `gorm:"size:2"`
	RegistrationDate mysql.NullTime
	Birthday         mysql.NullTime
	Gender           sql.NullString `gorm:"size:2"`  // e.g., "M", "F" ...
	Education        sql.NullString `gorm:"size:20"` // e.g., "High School"
	EnableEmail      sql.NullInt64  `gorm:"size:2"`
}

// OAuthAccount ...
type OAuthAccount struct {
	gorm.Model
	UserID    uint
	Type      string         `gorm:"size:10"`  // Facebook / Google ...
	AId       sql.NullString `gorm:"not null"` // user ID returned by OAuth services
	Email     sql.NullString `gorm:"size:100"`
	Name      sql.NullString `gorm:"size:80"`
	FirstName sql.NullString `gorm:"size:50"`
	LastName  sql.NullString `gorm:"size:50"`
	Gender    sql.NullString `gorm:"size:20"`
	Picture   sql.NullString // user profile photo url
	Birthday  sql.NullString
}

// ReporterAccount ...
type ReporterAccount struct {
	UserID        uint
	ID            uint   `gorm:"primary_key"`
	Email         string `gorm:"size:100;unique_index;not null"`
	Password      string `gorm:"not null"`
	Active        bool   `gorm:"default:false`
	ActivateToken string `gorm:"size:50"`
}
