package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	//"twreporter.org/go-api/globals"
	//log "github.com/Sirupsen/logrus"
)

// UsersBookmarks users and bookmarks many-to-many table
type UsersBookmarks struct {
	gorm.JoinTableHandler
	UserID     int
	BookmarkID int
	CreatedAt  time.Time
}

// Add - implement gorm.JoinTableHandlerInterface Add method
// Use gorm to create the record automatically populated with CreatedAt
func (*UsersBookmarks) Add(handler gorm.JoinTableHandlerInterface, db *gorm.DB, foreignValue interface{}, associationValue interface{}) error {
	foreignPrimaryKey, _ := strconv.Atoi(fmt.Sprint(db.NewScope(foreignValue).PrimaryKeyValue()))
	associationPrimaryKey, _ := strconv.Atoi(fmt.Sprint(db.NewScope(associationValue).PrimaryKeyValue()))

	return db.Create(&UsersBookmarks{
		UserID:     foreignPrimaryKey,
		BookmarkID: associationPrimaryKey,
		CreatedAt:  time.Now(),
	}).Error
}
