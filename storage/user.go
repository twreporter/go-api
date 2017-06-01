package storage

import (
	"database/sql"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
	"twreporter.org/go-api/configs/constants"
	"twreporter.org/go-api/models"
	"twreporter.org/go-api/utils"
)

// GetUserByID gets the user by its ID
func (s *GormMembershipStorage) GetUserByID(userID string) (models.User, error) {
	user := models.User{}

	// SELECT * FROM users WHERE ID = $userID
	err := s.db.First(&user, "id = ?", userID).Error
	return user, err
}

// GetOAuthData gets the corresponding OAuth by using the OAuth information
func (s *GormMembershipStorage) GetOAuthData(aid sql.NullString, aType string) (models.OAuthAccount, error) {
	log.Info("Getting the matching OAuth data", aid)
	oac := models.OAuthAccount{}
	err := s.db.Where(&models.OAuthAccount{Type: aType, AId: aid}).Last(&oac).Error
	if err != nil {
		log.Error("stroage.stroge_user.get_oauth_data.select_record_error:", err)
	}
	return oac, err
}

// GetUserDataByOAuth gets the corresponding user data by using the OAuth information
func (s *GormMembershipStorage) GetUserDataByOAuth(oac models.OAuthAccount) (models.User, error) {
	log.Info("Getting the matching User data")

	user := models.User{}

	matO, err := s.GetOAuthData(oac.AId, oac.Type)
	if err != nil {
		return user, err
	}

	err = s.db.Model(&matO).Related(&user).Error
	if err != nil {
		log.Error("stroage.storage_user.get_user_data_by_oauth.select_record_error: ", err)
	}
	return user, err
}

// GetReporterAccountData get the corresponding Reporter account by comparing email and password
func (s *GormMembershipStorage) GetReporterAccountData(email string) (*models.ReporterAccount, error) {
	log.WithFields(log.Fields{
		"email": email,
	}).Info("Getting the matching Reporter account data")

	ra := models.ReporterAccount{}
	err := s.db.Where(&models.ReporterAccount{Account: email}).Find(&ra).Error
	return &ra, err
}

// GetUserDataByReporterAccount get user data from user table by providing its reporter account data
func (s *GormMembershipStorage) GetUserDataByReporterAccount(ra *models.ReporterAccount) (*models.User, error) {
	log.Info("Getting the matching User data by reporter account")
	user := models.User{}
	err := s.db.Model(ra).Related(&user).Error
	return &user, err
}

// InsertUserByOAuth insert a new user into db after the oath loginin
func (s *GormMembershipStorage) InsertUserByOAuth(omodel models.OAuthAccount) models.User {
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

// InsertUserByReporterAccount insert a new user into db after the sign up
func (s *GormMembershipStorage) InsertUserByReporterAccount(raModel models.ReporterAccount) (models.User, error) {
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

// UpdateOAuthData updates the corresponding OAuth by using the OAuth information
func (s *GormMembershipStorage) UpdateOAuthData(newData models.OAuthAccount) (models.OAuthAccount, error) {
	log.Info("Getting the matching OAuth data", newData.AId)
	matO, err := s.GetOAuthData(newData.AId, newData.Type)
	if err != nil {
		return matO, err
	}
	matO.Email = newData.Email
	matO.Name = newData.Name
	matO.FirstName = newData.FirstName
	matO.LastName = newData.LastName
	matO.Gender = newData.Gender
	matO.Picture = newData.Picture
	s.db.Save(&matO)

	return matO, err
}

// UpdateReporterAccountPassword update password for a reporter account
func (s *GormMembershipStorage) UpdateReporterAccountPassword(ra *models.ReporterAccount, password string) (*models.ReporterAccount, error) {
	ra.Password = password
	err := s.db.Save(ra).Error
	return ra, err
}

// UpdateReporterAccountActive update password for a reporter account
func (s *GormMembershipStorage) UpdateReporterAccountActive(ra *models.ReporterAccount, active bool) (*models.ReporterAccount, error) {
	ra.Active = active
	err := s.db.Save(ra).Error
	return ra, err
}
