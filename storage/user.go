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
func (gs *GormStorage) GetUserByID(userID string) (models.User, error) {
	user := models.User{}

	// SELECT * FROM users WHERE ID = $userID
	if err := gs.db.First(&user, "id = ?", userID).Error; err != nil {
		return user, gs.NewStorageError(err, "GormStorage.GetUserByID", "storage.user.get_user_by_id.error")
	}

	return user, nil
}

// GetUserByEmail gets the user by its email
func (gs *GormStorage) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}

	// SELECT * FROM users WHERE email = $email
	err := gs.db.First(&user, "email = ?", email).Error
	if err != nil {
		return user, gs.NewStorageError(err, "GormStorage.GetOAuthData", "storage.user.get_o_auth_data.error")
	}
	return user, err
}

// GetOAuthData gets the corresponding OAuth by using the OAuth information
func (gs *GormStorage) GetOAuthData(aid sql.NullString, aType string) (models.OAuthAccount, error) {
	log.Info("Getting the matching OAuth data", aid)
	oac := models.OAuthAccount{}
	err := gs.db.Where(&models.OAuthAccount{Type: aType, AId: aid}).Last(&oac).Error
	if err != nil {
		return oac, gs.NewStorageError(err, "GormStorage.GetOAuthData", "storage.user.get_o_auth_data.error")
	}
	return oac, err
}

// GetUserDataByOAuth gets the corresponding user data by using the OAuth information
func (gs *GormStorage) GetUserDataByOAuth(oac models.OAuthAccount) (models.User, error) {
	log.Info("Getting the matching User data")

	user := models.User{}

	matO, err := gs.GetOAuthData(oac.AId, oac.Type)
	if err != nil {
		return user, err
	}

	err = gs.db.Model(&matO).Related(&user).Error
	if err != nil {
		log.Error("stroage.storage_user.get_user_data_by_oauth.select_record_error: ", err)
	}
	return user, err
}

// GetReporterAccountData get the corresponding Reporter account by comparing email and password
func (gs *GormStorage) GetReporterAccountData(email string) (models.ReporterAccount, error) {
	log.WithFields(log.Fields{
		"email": email,
	}).Info("Getting the matching Reporter account data")

	ra := models.ReporterAccount{}
	err := gs.db.Where(&models.ReporterAccount{Email: email}).Find(&ra).Error
	return ra, gs.NewStorageError(err, "GormStorage.GetReporterAccountData", "Getting account from reporter_accounts table occurs error")
}

// GetUserDataByReporterAccount get user data from user table by providing its reporter account data
func (gs *GormStorage) GetUserDataByReporterAccount(ra models.ReporterAccount) (models.User, error) {
	log.Info("Getting the matching User data by reporter account")
	user := models.User{}
	err := gs.db.Model(ra).Related(&user).Error
	return user, err
}

// InsertOAuthAccount insert  a new record into o_auth_accounts table
func (gs *GormStorage) InsertOAuthAccount(account models.OAuthAccount) error {
	err := gs.db.Create(&account).Error
	if err != nil {
		return gs.NewStorageError(err, "GormStorage.InsertOAuthAccount", "storage.user.create_oauth_account.error")
	}
	return nil
}

// InsertReporterAccount insert  a new record into reporter_accounts table
func (gs *GormStorage) InsertReporterAccount(account models.ReporterAccount) error {
	err := gs.db.Create(&account).Error
	if err != nil {
		return gs.NewStorageError(err, "GormStorage.InsertReporterAccount", "storage.user.create_reporter_account.error")
	}
	return nil
}

// InsertUserByOAuth insert a new user into db after the oath loginin
func (gs *GormStorage) InsertUserByOAuth(omodel models.OAuthAccount) models.User {
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
	gs.db.Create(&user)
	return user
}

// InsertUserByReporterAccount insert a new user into db after the sign up
func (gs *GormStorage) InsertUserByReporterAccount(raModel models.ReporterAccount) (models.User, error) {
	user := models.User{
		ReporterAccount:  raModel,
		Email:            utils.ToNullString(raModel.Email),
		RegistrationDate: mysql.NullTime{Time: time.Now(), Valid: true},
	}
	err := gs.db.Create(&user).Error
	return user, err
}

// UpdateOAuthData updates the corresponding OAuth by using the OAuth information
func (gs *GormStorage) UpdateOAuthData(newData models.OAuthAccount) (models.OAuthAccount, error) {
	log.Info("Getting the matching OAuth data", newData.AId)
	matO, err := gs.GetOAuthData(newData.AId, newData.Type)
	if err != nil {
		return matO, err
	}
	matO.Email = newData.Email
	matO.Name = newData.Name
	matO.FirstName = newData.FirstName
	matO.LastName = newData.LastName
	matO.Gender = newData.Gender
	matO.Picture = newData.Picture
	gs.db.Save(&matO)

	return matO, err
}

// UpdateReporterAccount update a reporter account
func (gs *GormStorage) UpdateReporterAccount(ra models.ReporterAccount) error {
	err := gs.db.Save(&ra).Error
	return err
}
