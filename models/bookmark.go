package models

import (
	"database/sql"
	"time"
)

// Bookmark this is bookmarks table description
type Bookmark struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at"`
	Href      string         `gorm:"size:512;unique_index;not null" json:"href"`
	Title     string         `gorm:"size:100;not null" json:"title"`
	Desc      sql.NullString `gorm:"size:250" json:"desc"`
	Thumbnail sql.NullString `gorm:"size:512" json:"thumbnail"`
}

// BookmarkForm this is the request BODY in form format
type BookmarkForm struct {
	Href      string `form:"href" binding:"required"`
	Title     string `form:"title" binding:"required"`
	Desc      string `form:"desc"`
	Thumbnail string `form:"thumbnail"`
}

// BookmarkJSON the is the request BODY in JSON format
type BookmarkJSON struct {
	Href      string `json:"href" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Desc      string `json:"desc"`
	Thumbnail string `json:"thumbnail"`
}
