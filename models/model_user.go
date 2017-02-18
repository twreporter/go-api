package models

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// User ...
type User struct {
	gorm.Model                      // contains fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	OAuthAccounts    []OAuthAccount `gorm:"ForeignKey:UserID"` // a user has multiple oauth accounts //
	Email            sql.NullString `gorm:"size:100"`
	FirstName        sql.NullString `gorm:"size:50"`
	LastName         sql.NullString `gorm:"size:50"`
	SecurityID       sql.NullString `gorm:"size:20"`
	PassportID       sql.NullString `gorm:"size:30"`
	City             sql.NullString `gorm:"size:45"`
	State            sql.NullString `gorm:"size:45"`
	Country          sql.NullString `gorm:"size:45"`
	Zip              sql.NullString `gorm:"size:20"`
	Address          sql.NullString
	Phone            sql.NullString `gorm:"size:20"`
	Privilege        int            `gorm:"size:3;not null"`
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
