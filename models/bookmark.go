package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

// Bookmark this is bookmakrs table description
type Bookmark struct {
	gorm.Model
	Href      string         `gorm:"size:512;unique_index;not null"`
	Title     string         `gorm:"size:100;not null"`
	Desc      sql.NullString `gorm:"size:250"`
	Thumbnail sql.NullString `gorm:"size:512"`
}
