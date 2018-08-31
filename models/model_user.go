package models

import (
	"database/sql"
	"encoding/json"
	"time"

	// log "github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
)

type NullString struct {
	sql.NullString
}

func NewNullString(s string) NullString {
	return NullString{
		sql.NullString{
			String: s,
			Valid:  s != "",
		},
	}
}

func (v *NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullString) UnmarshalJSON(data []byte) error {
	var s *string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s != nil {
		v.Valid = true
		v.String = *s
	} else {
		v.Valid = false
	}
	return nil
}

// User ...
type User struct {
	ID               uint            `gorm:"primary_key" json:"id"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
	DeletedAt        *time.Time      `json:"deleted_at"`
	OAuthAccounts    []OAuthAccount  `gorm:"ForeignKey:UserID"` // a user has multiple oauth accounts //
	ReporterAccount  ReporterAccount `gorm:"ForeignKey:UserID"`
	Bookmarks        []Bookmark      `gorm:"many2many:users_bookmarks;"`
	Email            NullString      `gorm:"size:100" json:"email"`
	FirstName        NullString      `gorm:"size:50" json:"firstname"`
	LastName         NullString      `gorm:"size:50" json:"lastname"`
	SecurityID       NullString      `gorm:"size:20" json:"security_id"`
	PassportID       NullString      `gorm:"size:30" json:"passport_id"`
	City             NullString      `gorm:"size:45" json:"city"`
	State            NullString      `gorm:"size:45" json:"state"`
	Country          NullString      `gorm:"size:45" json:"country"`
	Zip              NullString      `gorm:"size:20" json:"zip"`
	Address          NullString      `json:"address"`
	Phone            NullString      `gorm:"size:20" json:"phone"`
	Privilege        int             `gorm:"type:int(5);not null" json:"privilege"`
	RegistrationDate mysql.NullTime  `json:"registration_date"`
	Birthday         mysql.NullTime  `json:"birthday"`
	Gender           NullString      `gorm:"size:2" json:"gender"`     // e.g., "M", "F" ...
	Education        NullString      `gorm:"size:20" json:"education"` // e.g., "High School"
	EnableEmail      int             `gorm:"type:int(5);size:2" json:"enable_email"`
}

// OAuthAccount ...
type OAuthAccount struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UserID    uint       `json:"user_id"`
	Type      string     `gorm:"size:10" json:"type"`  // Facebook / Google ...
	AId       NullString `gorm:"not null" json:"a_id"` // user ID returned by OAuth services
	Email     NullString `gorm:"size:100" json:"email"`
	Name      NullString `gorm:"size:80" json:"name"`
	FirstName NullString `gorm:"size:50" json:"firstname"`
	LastName  NullString `gorm:"size:50" json:"lastname"`
	Gender    NullString `gorm:"size:20" json:"gender"`
	Picture   NullString `json:"picture"` // user profile photo url
	Birthday  NullString `json:"birthday"`
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
