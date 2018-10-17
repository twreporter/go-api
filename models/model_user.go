package models

import (
	"time"

	// log "github.com/Sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
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
	Email            null.String     `gorm:"size:100" json:"email"`
	FirstName        null.String     `gorm:"size:50" json:"firstname"`
	LastName         null.String     `gorm:"size:50" json:"lastname"`
	SecurityID       null.String     `gorm:"size:20" json:"security_id"`
	PassportID       null.String     `gorm:"size:30" json:"passport_id"`
	City             null.String     `gorm:"size:45" json:"city"`
	State            null.String     `gorm:"size:45" json:"state"`
	Country          null.String     `gorm:"size:45" json:"country"`
	Zip              null.String     `gorm:"size:20" json:"zip"`
	Address          null.String     `json:"address"`
	Phone            null.String     `gorm:"size:20" json:"phone"`
	Privilege        int             `gorm:"type:int(5);not null" json:"privilege"`
	RegistrationDate null.Time       `json:"registration_date"`
	Birthday         null.Time       `json:"birthday"`
	Gender           null.String     `gorm:"size:2" json:"gender"`     // e.g., "M", "F" ...
	Education        null.String     `gorm:"size:20" json:"education"` // e.g., "High School"
	EnableEmail      int             `gorm:"type:int(5);size:2" json:"enable_email"`
}

// OAuthAccount ...
type OAuthAccount struct {
	ID        uint        `gorm:"primary_key" json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt *time.Time  `json:"deleted_at"`
	UserID    uint        `json:"user_id"`
	Type      string      `gorm:"size:10" json:"type"`  // Facebook / Google ...
	AId       null.String `gorm:"not null" json:"a_id"` // user ID returned by OAuth services
	Email     null.String `gorm:"size:100" json:"email"`
	Name      null.String `gorm:"size:80" json:"name"`
	FirstName null.String `gorm:"size:50" json:"firstname"`
	LastName  null.String `gorm:"size:50" json:"lastname"`
	Gender    null.String `gorm:"size:20" json:"gender"`
	Picture   null.String `json:"picture"` // user profile photo url
	Birthday  null.String `json:"birthday"`
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
