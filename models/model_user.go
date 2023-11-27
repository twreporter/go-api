package models

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

// User ...
type User struct {
	ID                  uint              `gorm:"primary_key" json:"id"`
	CreatedAt           time.Time         `json:"created_at"`
	UpdatedAt           time.Time         `json:"updated_at"`
	DeletedAt           *time.Time        `json:"deleted_at"`
	OAuthAccounts       []OAuthAccount    `gorm:"ForeignKey:UserID"` // a user has multiple oauth accounts //
	ReporterAccount     ReporterAccount   `gorm:"ForeignKey:UserID"`
	Bookmarks           []Bookmark        `gorm:"many2many:users_bookmarks;"`
	MailGroups          []UsersMailgroups `gorm:"one2many:users_mailgroups;"`
	Email               null.String       `gorm:"size:100" json:"email"`
	FirstName           null.String       `gorm:"size:50" json:"firstname"`
	LastName            null.String       `gorm:"size:50" json:"lastname"`
	Nickname            null.String       `gorm:"size:50" json:"nickname"`
	SecurityID          null.String       `gorm:"size:20" json:"security_id"`
	PassportID          null.String       `gorm:"size:30" json:"passport_id"`
	Title               null.String       `gorm:"size:30" json:"title"`
	LegalName           null.String       `gorm:"size:50" json:"legal_name"`
	City                null.String       `gorm:"size:45" json:"city"`
	State               null.String       `gorm:"size:45" json:"state"`
	Country             null.String       `gorm:"size:45" json:"country"`
	Zip                 null.String       `gorm:"size:20" json:"zip"`
	Address             null.String       `json:"address"`
	Phone               null.String       `gorm:"size:20" json:"phone"`
	Privilege           int               `gorm:"type:int(5);not null" json:"privilege"`
	RegistrationDate    null.Time         `json:"registration_date"`
	Birthday            null.Time         `json:"birthday"`
	Gender              null.String       `gorm:"size:2" json:"gender"` // e.g., "M", "F", "X, "U"
	AgeRange            null.String       `gorm:"type:ENUM('less_than_18', '18_to_24', '25_to_34', '35_to_44', '45_to_54', '55_to_64', 'above_65')" json:"age_range"`
	Education           null.String       `gorm:"size:20" json:"education"` // e.g., "High School"
	EnableEmail         int               `gorm:"type:int(5);size:2" json:"enable_email"`
	ReadPreference      null.String       `gorm:"type:SET('international', 'cross_straits', 'human_right', 'society', 'environment', 'education', 'politics', 'economy', 'culture', 'art', 'life', 'health', 'sport', 'all')" json:"read_preference"` // e.g. "international, art, sport"
	WordsForTwreporter  null.String       `gorm:"size:255" json:"words_for_twreporter"`
	Roles               []Role            `gorm:"many2many:users_roles" json:"roles"`
	Activated           null.Time         `json:"activated"`
	Source              null.String       `gorm:"type:SET('ntch')" json:"source"`
	AgreeDataCollection bool              `gorm:"type:tinyint(1);default:1" json:"agree_data_collection"`
	ReadPostsCount      int               `gorm:"type:int(10);unsigned" json:"read_posts_count"`
	ReadPostsSec        int               `gorm:"type:int(10);unsigned" json:"read_posts_sec"`
}

// Role represents a user role
type Role struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	NameEn    string    `json:"name_en"`
	Key       string    `json:"key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Weight    int       `json:"weight"`
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
