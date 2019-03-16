package storage

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"

	"gopkg.in/guregu/null.v3"
	"twreporter.org/go-api/configs/constants"
	"twreporter.org/go-api/models"
)

// GetUserByID gets the user by its ID
func (gs *GormStorage) GetUserByID(userID string) (models.User, error) {
	user := models.User{}

	// SELECT * FROM users WHERE ID = $userID
	if err := gs.db.First(&user, "id = ?", userID).Error; err != nil {
		return user, gs.NewStorageError(err, "GormStorage.GetUserByID", fmt.Sprintf("get user(id: %s) error", userID))
	}

	return user, nil
}

// GetUserByEmail gets the user by its email
func (gs *GormStorage) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}

	// SELECT * FROM users WHERE email = $email
	err := gs.db.First(&user, "email = ?", email).Error
	if err != nil {
		return user, gs.NewStorageError(err, "GormStorage.GetOAuthData", fmt.Sprintf("get user(email: %s) error", email))
	}
	return user, err
}

// GetOAuthData gets the corresponding OAuth by using the OAuth information
func (gs *GormStorage) GetOAuthData(aid null.String, aType string) (models.OAuthAccount, error) {
	log.Info("Getting the matching OAuth data", aid)
	oac := models.OAuthAccount{}
	err := gs.db.Where(&models.OAuthAccount{Type: aType, AId: aid}).Last(&oac).Error
	if err != nil {
		return oac, gs.NewStorageError(err, "GormStorage.GetOAuthData", "get oauth account error")
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
	return ra, gs.NewStorageError(err, "GormStorage.GetReporterAccountData", fmt.Sprintf("get reporter account(email: %s) error", email))
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
		return gs.NewStorageError(err, "GormStorage.InsertOAuthAccount", fmt.Sprint("create oauth account error"))
	}
	return nil
}

// InsertReporterAccount insert  a new record into reporter_accounts table
func (gs *GormStorage) InsertReporterAccount(account models.ReporterAccount) error {
	err := gs.db.Create(&account).Error
	if err != nil {
		return gs.NewStorageError(err, "GormStorage.InsertReporterAccount", fmt.Sprint("create reporter account error"))
	}
	return nil
}

// InsertUserByOAuth insert a new user into db after the oath loginin
func (gs *GormStorage) InsertUserByOAuth(omodel models.OAuthAccount) (user models.User, err error) {
	log.Info("Inserting user data")
	user = models.User{
		OAuthAccounts:    []models.OAuthAccount{omodel},
		Email:            omodel.Email,
		FirstName:        omodel.FirstName,
		LastName:         omodel.LastName,
		Gender:           omodel.Gender,
		Privilege:        constants.PrivilegeRegistered,
		RegistrationDate: null.TimeFrom(time.Now()),
	}
	err = gs.db.Create(&user).Error
	return user, err
}

// InsertUserByReporterAccount insert a new user into db after the sign up
func (gs *GormStorage) InsertUserByReporterAccount(raModel models.ReporterAccount) (models.User, error) {
	user := models.User{
		ReporterAccount:  raModel,
		Email:            null.StringFrom(raModel.Email),
		RegistrationDate: null.NewTime(time.Now(), true),
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
	err := gs.db.Model(&ra).Updates(&ra).Error
	return err
}
