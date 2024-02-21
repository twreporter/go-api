package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

// UsersBookmarks users and bookmarks many-to-many table
type UsersBookmarks struct {
	gorm.JoinTableHandler
	UserID     int
	BookmarkID int
	CreatedAt  time.Time
	PostID     string
}

// Add - implement gorm.JoinTableHandlerInterface Add method
// Use gorm to create the record automatically populated with CreatedAt
func (*UsersBookmarks) Add(handler gorm.JoinTableHandlerInterface, db *gorm.DB, foreignValue interface{}, associationValue interface{}) error {
	foreignPrimaryKey, _ := strconv.Atoi(fmt.Sprint(db.NewScope(foreignValue).PrimaryKeyValue()))
	associationPrimaryKey, _ := strconv.Atoi(fmt.Sprint(db.NewScope(associationValue).PrimaryKeyValue()))

	postIDField, ok := db.NewScope(associationValue).FieldByName("post_id")
	var postID string
	if ok != false {
		postID = postIDField.Field.String()
	}

	return db.Create(&UsersBookmarks{
		UserID:     foreignPrimaryKey,
		BookmarkID: associationPrimaryKey,
		CreatedAt:  time.Now(),
		PostID:     postID,
	}).Error
}
