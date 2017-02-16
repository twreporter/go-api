package storage

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/models"
)

// NewUserStorage initializes the user storage
func NewUserStorage(db *gorm.DB) *UserStorage {
	return &UserStorage{db}
}

// UserStorage stores all users
type UserStorage struct {
	db *gorm.DB
}

// // InsertUserByOAuth insert a new user into db after the oath login
// func (s UserStorage) InsertUserByOAuth(id string) (*User, error) {
//
// }

// GetUserByOAuth gets the corresponding user by using the oauth information
func (s UserStorage) GetUserByOAuth(aid string) models.User {
	log.Info("Getting the matching user data")
	user := models.User{}
	oac := models.OAuthAccount{Type: "Facebook", AId: aid}
	s.db.Model(&user).Related(&oac)
	return user
}
