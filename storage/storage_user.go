package storage

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/configs/constants"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/utils"
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
func (s UserStorage) GetOAuthData(aid sql.NullString) models.OAuthAccount {
	log.Info("Getting the matching OAuth data", aid)
	oac := models.OAuthAccount{}
	s.db.Where(&models.OAuthAccount{Type: constants.Facebook, AId: aid}).First(&oac)
	return oac
}

// GetUserDataByOAuth gets the corresponding user data by using the OAuth information
func (s UserStorage) GetUserDataByOAuth(oac models.OAuthAccount) models.User {
	log.Info("Getting the matching User data")
	matO := s.GetOAuthData(oac.AId)
	user := models.User{}
	s.db.Model(&matO).Related(&user)
	return user
}

// UpdateOAuthData updates the corresponding OAuth by using the OAuth information
func (s UserStorage) UpdateOAuthData(newData models.OAuthAccount) models.OAuthAccount {
	log.Info("Getting the matching OAuth data", newData.AId)
	matO := s.GetOAuthData(newData.AId)
	matO.Email = newData.Email
	matO.Name = newData.Name
	matO.FirstName = newData.FirstName
	matO.LastName = newData.LastName
	matO.Gender = newData.Gender
	matO.Picture = newData.Picture
	s.db.Save(&matO)
	return matO
}

// InsertUserByReporterAccount insert a new user into db after the sign up
func (s UserStorage) InsertUserByReporterAccount(raModel models.ReporterAccount) (models.User, error) {
	log.WithFields(log.Fields{
		"account":       raModel.Email,
		"password":      raModel.Password,
		"Active":        raModel.Active,
		"ActivateToken": raModel.ActivateToken,
	}).Info("Inserting user data")
	user := models.User{
		ReporterAccount:  raModel,
		Email:            utils.ToNullString(raModel.Email),
		RegistrationDate: mysql.NullTime{Time: time.Now(), Valid: true},
	}
	err := s.db.Create(&user).Error
	return user, err
}

// GetReporterAccountData get the corresponding Reporter account by comparing email and password
func (s UserStorage) GetReporterAccountData(email string) (models.ReporterAccount, error) {
	log.WithFields(log.Fields{
		"email": email,
	}).Info("Getting the matching Reporter account data")

	rc := models.ReporterAccount{}
	err := s.db.Where(&models.ReporterAccount{Email: email}).Find(&rc).Error
	return rc, err
}
