package models

import (
	"time"
)

// UsersPostsReadingCount: users and reading posts count are one-to-many table
type UsersPostsReadingCount struct {
	UserID     int
	PostID     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
}

// UsersPostsReadingTime: users and reading posts time are one-to-many table
type UsersPostsReadingTime struct {
	UserID     int
	PostID     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time
	Seconds    int
}
