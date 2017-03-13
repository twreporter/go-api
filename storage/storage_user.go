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

// UserStorage this is an interface defines methods for users, reporter_accounts and o_auth_accounts tables
type UserStorage interface {
	InsertUserByOAuth(models.OAuthAccount) models.User
	GetOAuthData(sql.NullString, string) models.OAuthAccount
	GetUserDataByOAuth(models.OAuthAccount) (models.User, error)
	UpdateOAuthData(models.OAuthAccount) models.OAuthAccount
	InsertUserByReporterAccount(models.ReporterAccount) (models.User, error)
	GetReporterAccountData(string) (*models.ReporterAccount, error)
	GetUserDataByReporterAccount(*models.ReporterAccount) (*models.User, error)
	UpdateReporterAccountPassword(*models.ReporterAccount, string) (*models.ReporterAccount, error)
	UpdateReporterAccountActive(*models.ReporterAccount, bool) (*models.ReporterAccount, error)
}

// NewGormUserStorage this initializes the user storage
func NewGormUserStorage(db *gorm.DB) UserStorage {
	return &gormUserStorage{db}
}

// gormUserStorage this implements UserStorage interface
type gormUserStorage struct {
	db *gorm.DB
}

// InsertUserByOAuth insert a new user into db after the oath loginin
func (s *gormUserStorage) InsertUserByOAuth(omodel models.OAuthAccount) models.User {
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
func (s *gormUserStorage) GetOAuthData(aid sql.NullString, aType string) models.OAuthAccount {
	log.Info("Getting the matching OAuth data", aid)
	oac := models.OAuthAccount{}
	s.db.Where(&models.OAuthAccount{Type: aType, AId: aid}).First(&oac)
	return oac
}

// GetUserDataByOAuth gets the corresponding user data by using the OAuth information
func (s *gormUserStorage) GetUserDataByOAuth(oac models.OAuthAccount) (models.User, error) {
	log.Info("Getting the matching User data")
	matO := s.GetOAuthData(oac.AId, oac.Type)
	user := models.User{}
	err := s.db.Model(&matO).Related(&user).Error
	return user, err
}

// UpdateOAuthData updates the corresponding OAuth by using the OAuth information
func (s *gormUserStorage) UpdateOAuthData(newData models.OAuthAccount) models.OAuthAccount {
	log.Info("Getting the matching OAuth data", newData.AId)
	matO := s.GetOAuthData(newData.AId, newData.Type)
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
func (s *gormUserStorage) InsertUserByReporterAccount(raModel models.ReporterAccount) (models.User, error) {
	log.WithFields(log.Fields{
		"account":       raModel.Account,
		"password":      raModel.Password,
		"Active":        raModel.Active,
		"ActivateToken": raModel.ActivateToken,
	}).Info("Inserting user data")
	user := models.User{
		ReporterAccount:  raModel,
		Email:            utils.ToNullString(raModel.Account),
		RegistrationDate: mysql.NullTime{Time: time.Now(), Valid: true},
	}
	err := s.db.Create(&user).Error
	return user, err
}

// GetReporterAccountData get the corresponding Reporter account by comparing email and password
func (s *gormUserStorage) GetReporterAccountData(email string) (*models.ReporterAccount, error) {
	log.WithFields(log.Fields{
		"email": email,
	}).Info("Getting the matching Reporter account data")

	ra := models.ReporterAccount{}
	err := s.db.Where(&models.ReporterAccount{Account: email}).Find(&ra).Error
	return &ra, err
}

// GetUserDataByReporterAccount get user data from user table by providing its reporter account data
func (s *gormUserStorage) GetUserDataByReporterAccount(ra *models.ReporterAccount) (*models.User, error) {
	log.Info("Getting the matching User data by reporter account")
	user := models.User{}
	err := s.db.Model(ra).Related(&user).Error
	return &user, err
}

// UpdateReporterAccountPassword update password for a reporter account
func (s *gormUserStorage) UpdateReporterAccountPassword(ra *models.ReporterAccount, password string) (*models.ReporterAccount, error) {
	ra.Password = password
	err := s.db.Save(ra).Error
	return ra, err
}

// UpdateReporterAccountActive update password for a reporter account
func (s *gormUserStorage) UpdateReporterAccountActive(ra *models.ReporterAccount, active bool) (*models.ReporterAccount, error) {
	ra.Active = active
	err := s.db.Save(ra).Error
	return ra, err
}
