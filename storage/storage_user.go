package storage

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/configs/constants"
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

// InsertUserByOAuth insert a new user into db after the oath loginin
func (s UserStorage) InsertUserByOAuth(omodel models.OAuthAccount) models.User {
	log.Info("Inserting user data")
	user := models.User{
		OAuthAccounts:    []models.OAuthAccount{omodel},
		Email:            omodel.Email,
		FirstName:        omodel.FirstName,
		LastName:         omodel.LastName,
		Gender:           omodel.Gender,
		Privilege:        constants.PrivilegeRegistered,
		RegistrationDate: mysql.NullTime{Time: time.Now(), Valid: true},
	}
	s.db.Create(&user)
	return user
}

// GetOAuthData gets the corresponding OAuth by using the OAuth information
func (s UserStorage) GetOAuthData(aid string) models.OAuthAccount {
	log.Info("Getting the matching OAuth data")
	oac := models.OAuthAccount{}
	s.db.Where(&models.OAuthAccount{Type: constants.Facebook}).First(&oac)
	return oac
}

// GetUserDataByOAuth gets the corresponding user data by using the OAuth information
func (s UserStorage) GetUserDataByOAuth(oac models.OAuthAccount) models.User {
	log.Info("Getting the matching User data")
	user := models.User{}
	s.db.Model(&oac).Related(&user)
	return user
}
