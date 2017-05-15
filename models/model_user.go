package models

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// User ...
type User struct {
	ID               uint `gorm:"primary_key"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
	OAuthAccounts    []OAuthAccount  `gorm:"ForeignKey:UserID"` // a user has multiple oauth accounts //
	ReporterAccount  ReporterAccount `gorm:"ForeignKey:UserID"`
	Bookmarks        []Bookmark      `gorm:"many2many:users_bookmarks;"`
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
	Privilege        int            `gorm:"type:int(5);not null"`
	RegistrationDate mysql.NullTime
	Birthday         mysql.NullTime
	Gender           sql.NullString `gorm:"size:2"`  // e.g., "M", "F" ...
	Education        sql.NullString `gorm:"size:20"` // e.g., "High School"
	EnableEmail      int            `gorm:"type:int(5);size:2"`
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
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
	ID            uint   `gorm:"primary_key"`
	Account       string `gorm:"size:100;unique_index;not null"`
	Password      string `gorm:"not null"`
	Active        bool   `gorm:"default:0"`
	ActivateToken string `gorm:"size:50"`
}
