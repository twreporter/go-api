package models

import (
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
)

// User ...
type User struct {
	ID               uint            `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        *time.Time      `json:"deleted_at"`
	OAuthAccounts    []OAuthAccount  `gorm:"ForeignKey:UserID"` // a user has multiple oauth accounts //
	ReporterAccount  ReporterAccount `gorm:"ForeignKey:UserID"`
	Bookmarks        []Bookmark      `gorm:"many2many:users_bookmarks;"`
	Email            sql.NullString  `gorm:"size:100" json:"email"`
	FirstName        sql.NullString  `gorm:"size:50" json:"firstname"`
	LastName         sql.NullString  `gorm:"size:50" json:"lastname"`
	SecurityID       sql.NullString  `gorm:"size:20" json:"security_id"`
	PassportID       sql.NullString  `gorm:"size:30" json:"passport_id"`
	City             sql.NullString  `gorm:"size:45" json:"city"`
	State            sql.NullString  `gorm:"size:45" json:"state"`
	Country          sql.NullString  `gorm:"size:45" json:"country"`
	Zip              sql.NullString  `gorm:"size:20" json:"zip"`
	Address          sql.NullString  `json:"address"`
	Phone            sql.NullString  `gorm:"size:20" json:"phone"`
	Privilege        int             `gorm:"type:int(5);not null" json:"privilege"`
	RegistrationDate mysql.NullTime  `json:"registration_date"`
	Birthday         mysql.NullTime  `json:"birthday"`
	Gender           sql.NullString  `gorm:"size:2" json:"gender"`     // e.g., "M", "F" ...
	Education        sql.NullString  `gorm:"size:20" json:"education"` // e.g., "High School"
	EnableEmail      int             `gorm:"type:int(5);size:2" json:"enable_email"`
}

// OAuthAccount ...
type OAuthAccount struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at"`
	UserID    uint           `json:"user_id"`
	Type      string         `gorm:"size:10" json:"type"`  // Facebook / Google ...
	AId       sql.NullString `gorm:"not null" json:"a_id"` // user ID returned by OAuth services
	Email     sql.NullString `gorm:"size:100" json:"email"`
	Name      sql.NullString `gorm:"size:80" json:"name"`
	FirstName sql.NullString `gorm:"size:50" json:"firstname"`
	LastName  sql.NullString `gorm:"size:50" json:"lastname"`
	Gender    sql.NullString `gorm:"size:20" json:"gender"`
	Picture   sql.NullString `json:"picture"` // user profile photo url
	Birthday  sql.NullString `json:"birthday"`
}

// ReporterAccount ...
type ReporterAccount struct {
	UserID        uint       `json:"user_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
	ID            uint       `gorm:"primary_key" json:"id"`
	Email         string     `gorm:"size:100;unique_index;not null" json:"email"`
	ActivateToken string     `gorm:"size:50" json:"activate_token"`
	ActExpTime    time.Time  `json:"-"`
}
