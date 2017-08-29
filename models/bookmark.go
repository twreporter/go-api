package models

import (
	//"database/sql"
	"time"
)

// Bookmark this is bookmarks table description
type Bookmark struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Slug      string     `gorm:"size:100;unique_index;not null" json:"slug" form:"slug"`
	Title     string     `gorm:"size:100;not null" json:"title" from:"title"`
	Desc      string     `gorm:"size:250" json:"desc" form:"desc"`
	Style     string     `gorm:"size:100" json:"style" form:"style"`
	External  bool       `gorm:"default:0" json:"external" form:"external"`
	Thumbnail string     `gorm:"size:512" json:"thumbnail" form:"thumbnail"`
}
