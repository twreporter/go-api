package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
)

// Bookmark this is bookmarks table description
type Bookmark struct {
	gorm.Model
	Href      string         `gorm:"size:512;unique_index;not null"`
	Title     string         `gorm:"size:100;not null"`
	Desc      sql.NullString `gorm:"size:250"`
	Thumbnail sql.NullString `gorm:"size:512"`
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
