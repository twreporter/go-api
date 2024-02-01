package models

import (
	//"database/sql"
	"time"
)

// Bookmark this is bookmarks table description
type Bookmark struct {
	ID         uint       `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	Users      []User     `gorm:"many2many:users_bookmarks;"`
	Slug       string     `gorm:"size:100;not null" json:"slug" form:"slug" binding:"required"`
	Title      string     `gorm:"size:100;not null" json:"title" form:"title" binding:"required"`
	Desc       string     `gorm:"size:250" json:"desc" form:"desc"`
	Host       string     `gorm:"size:100;not null;" json:"host" form:"host" binding:"required"`
	Category   string     `gorm:"size:20" json:"category" form:"category"`
	IsExternal bool       `gorm:"default:0" json:"is_external" form:"is_external"`
	Thumbnail  string     `gorm:"size:512" json:"thumbnail" form:"thumbnail" binding:"required"`
	Authors    string     `gorm:"size:250" json:"authors" form:"authors"`
	PubDate    uint       `gorm:"not null;default:0" json:"published_date" form:"published_date"`
	PostID     string     `gorm:"size:50;not null" json:"post_id" form:"post_id"`
}

type UserBookmark struct {
	AddedAt    time.Time `json:"added_at" db:"users_bookmarks.added_at"`
	PostID     string    `json:"post_id" db:"users_bookmarks.post_id"`
	ID         uint      `json:"id" db:"bookmarks.id"`
	Slug       string    `json:"slug" db:"bookmarks.slug"`
	Title      string    `json:"title" db:"bookmarks.title"`
	Desc       string    `json:"desc" db:"bookmarks.desc"`
	Host       string    `json:"host" db:"bookmarks.host"`
	Category   string    `json:"category" db:"bookmarks.category"`
	IsExternal bool      `json:"is_external" db:"bookmarks.is_external"`
	Thumbnail  string    `json:"thumbnail" db:"bookmarks.thumbnail"`
	Authors    string    `json:"authors" db:"bookmarks.authors"`
	PubDate    uint      `json:"published_date" db:"bookmarks.published_date"`
}
