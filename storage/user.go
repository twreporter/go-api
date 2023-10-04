package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"

	"github.com/twreporter/go-api/configs/constants"
	"github.com/twreporter/go-api/models"
)

// GetUserByID gets the user by its ID
func (gs *GormStorage) GetUserByID(userID string) (models.User, error) {
	user := models.User{}

	// SELECT * FROM users WHERE ID = $userID and join roles and user_mailgroups tables
	err := gs.db.Preload("Roles").Preload("MailGroups").First(&user, "id = ?", userID).Error
	if err != nil {
		return user, errors.Wrap(err, fmt.Sprintf("get user(id: %s) error", userID))
	}

	return user, nil
}

// GetUserByEmail gets the user by its email
func (gs *GormStorage) GetUserByEmail(email string) (models.User, error) {
	user := models.User{}

	// SELECT * FROM users WHERE email = $email
	err := gs.db.First(&user, "email = ?", email).Error
	if err != nil {
		return user, errors.Wrap(err, fmt.Sprintf("get user(email: %s) error", email))
	}

	return user, nil
}

// GetOAuthData gets the corresponding OAuth by using the OAuth information
func (gs *GormStorage) GetOAuthData(aid null.String, aType string) (models.OAuthAccount, error) {
	log.Debug("Getting the matching OAuth data", aid)
	oac := models.OAuthAccount{}
	err := gs.db.Where(&models.OAuthAccount{Type: aType, AId: aid}).Last(&oac).Error
	if err != nil {
		return oac, errors.Wrap(err, "get oauth account error")
	}

	return oac, nil
}

// GetUserDataByOAuth gets the corresponding user data by using the OAuth information
func (gs *GormStorage) GetUserDataByOAuth(oac models.OAuthAccount) (models.User, error) {
	log.Debug("Getting the matching User data")

	user := models.User{}

	matO, err := gs.GetOAuthData(oac.AId, oac.Type)
	if err != nil {
		return user, err
	}

	err = gs.db.Model(&matO).Related(&user).Error
	if err != nil {
		return user, errors.WithStack(err)
	}

	return user, nil
}

// GetReporterAccountData get the corresponding Reporter account by comparing email and password
func (gs *GormStorage) GetReporterAccountData(email string) (models.ReporterAccount, error) {
	log.WithFields(log.Fields{
		"email": email,
	}).Debug("Getting the matching Reporter account data")

	ra := models.ReporterAccount{}
	err := gs.db.Where(&models.ReporterAccount{Email: email}).Find(&ra).Error

	if err != nil {
		return ra, errors.Wrap(err, fmt.Sprintf("get reporter account(email: %s) error", email))
	}

	return ra, nil
}

// GetUserDataByReporterAccount get user data from user table by providing its reporter account data
func (gs *GormStorage) GetUserDataByReporterAccount(ra models.ReporterAccount) (models.User, error) {
	log.Debug("Getting the matching User data by reporter account")
	user := models.User{}
	err := gs.db.Model(ra).Related(&user).Error
	if err != nil {
		return user, errors.WithStack(err)
	}

	return user, nil
}

// InsertOAuthAccount insert  a new record into o_auth_accounts table
func (gs *GormStorage) InsertOAuthAccount(account models.OAuthAccount) error {
	err := gs.db.Create(&account).Error
	if err != nil {
		return errors.Wrap(err, "create oauth account error")
	}
	return nil
}

// InsertReporterAccount insert  a new record into reporter_accounts table
func (gs *GormStorage) InsertReporterAccount(account models.ReporterAccount) error {
	err := gs.db.Create(&account).Error
	if err != nil {
		return errors.Wrap(err, "create reporter account error")
	}
	return nil
}

// InsertUserByOAuth insert a new user into db after the oath loginin
func (gs *GormStorage) InsertUserByOAuth(omodel models.OAuthAccount) (user models.User, err error) {
	log.Debug("Inserting user data")
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
	if err != nil {
		return user, errors.WithStack(err)
	}

	// Call AssignRoleToUser to assign role to user
	err = gs.AssignRoleToUser(user, constants.RoleExplorer)
	if err != nil {
		return user, errors.WithStack(err)
	}

	return user, nil
}

// InsertUserByReporterAccount insert a new user into db after the sign up
func (gs *GormStorage) InsertUserByReporterAccount(raModel models.ReporterAccount) (models.User, error) {
	user := models.User{
		ReporterAccount:  raModel,
		Email:            null.StringFrom(raModel.Email),
		RegistrationDate: null.NewTime(time.Now(), true),
	}
	err := gs.db.Create(&user).Error

	if err != nil {
		return user, errors.WithStack(err)
	}

	// Call AssignRoleToUser to assign role to user
	err = gs.AssignRoleToUser(user, constants.RoleExplorer)
	if err != nil {
		return user, errors.WithStack(err)
	}

	return user, nil
}

// UpdateOAuthData updates the corresponding OAuth by using the OAuth information
func (gs *GormStorage) UpdateOAuthData(newData models.OAuthAccount) (models.OAuthAccount, error) {
	log.Debug("Getting the matching OAuth data", newData.AId)
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
	err = gs.db.Save(&matO).Error

	if err != nil {
		return matO, errors.WithStack(err)
	}

	return matO, nil
}

// UpdateReporterAccount update a reporter account
func (gs *GormStorage) UpdateReporterAccount(ra models.ReporterAccount) error {
	err := gs.db.Model(&ra).Updates(&ra).Error

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// UpdateUser update a user
func (gs *GormStorage) UpdateUser(user models.User) error {
	err := gs.db.Model(&user).Updates(&user).Error

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// UpdateReadPreferenceOfUser function will update the read preference of a user
func (gs *GormStorage) UpdateReadPreferenceOfUser(userID string, readPreference []string) error {
	tx := gs.db.Begin() // Start the transaction

	// Check if the user exists
	var count int64
	if err := tx.Model(&models.User{}).Where("id = ?", userID).Count(&count).Error; err != nil {
		tx.Rollback() // Rollback the transaction if an error occurs
		return errors.Wrap(err, fmt.Sprintf("failed to check user existence (id: %s)", userID))
	}

	if count == 0 {
		tx.Rollback() // Rollback the transaction if the user doesn't exist
		return fmt.Errorf("user with ID %s does not exist", userID)
	}

	// Update the user's read preference
	if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("read_preference", null.StringFrom(strings.Join(readPreference, ","))).Error; err != nil {
		tx.Rollback() // Rollback the transaction if an error occurs
		return errors.Wrap(err, fmt.Sprintf("failed to update user's read preference (id: %s)", userID))
	}

	// Update the user's activated time to now
	if err := tx.Model(&models.User{}).Where("id = ?", userID).Update("activated", time.Now()).Error; err != nil {
		tx.Rollback() // Rollback the transaction if an error occurs
		return errors.Wrap(err, fmt.Sprintf("failed to update user's activated time (id: %s)", userID))
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback() // Rollback the transaction if an error occurs during commit
		return errors.Wrap(err, "failed to commit transaction")
	}

	return nil
}

// AssignRoleToUser assigns a role to a user, removing all existing roles
func (gs *GormStorage) AssignRoleToUser(user models.User, roleKey string) error {
	// Check if the role exists
	var role models.Role
	if err := gs.db.Where("`key` = ?", roleKey).First(&role).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to find role with Key: %s", roleKey))
	}

	// Replace existing roles with the new role
	association := gs.db.Model(&user).Association("Roles")
	association.Replace(&role)

	if association.Error != nil {
		return errors.Wrap(association.Error, "failed to assign role to user")
	}

	return nil
}

// GetRole retrieves the role of a user from the database based on the user's ID or unique identifier.
func (gs *GormStorage) GetRole(user models.User) (models.Role, error) {
	// Initialize an empty role variable to store the fetched role.
	var role models.Role

	// Get the associated roles for the user.
	association := gs.db.Model(&user).Association("Roles")

	if association.Error != nil {
		return role, errors.Wrap(association.Error, fmt.Sprintf("failed to get roles for user with email: %s", user.Email.String))
	}

	// Check if the user has any roles.
	if association.Count() == 0 {
		return role, errors.New(fmt.Sprintf("user with email: %s does not have any roles", user.Email.String))
	}

	// Get the first role associated with the user.
	if err := gs.db.Model(&user).Association("Roles").Find(&role).Error; err != nil {
		return role, errors.Wrap(err, fmt.Sprintf("failed to find role for user with email: %s", user.Email.String))
	}

	return role, nil
}

// HasRole checks if a specific role exists for a user based on the user's ID and the role key.
func (gs *GormStorage) HasRole(user models.User, roleKey string) (bool, error) {
	// Check if the role exists.
	var role models.Role
	if err := gs.db.Where("`key` = ?", roleKey).First(&role).Error; err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("failed to find role with Key: %s", roleKey))
	}

	// Fetch all the roles associated with the user.
	var roles []models.Role
	if err := gs.db.Model(&user).Association("Roles").Find(&roles).Error; err != nil {
		return false, errors.Wrap(err, fmt.Sprintf("failed to get roles for user with email: %s", user.Email.String))
	}

	// Check if the role with the given key exists in the retrieved roles.
	for _, r := range roles {
		if r.Key == roleKey {
			return true, nil
		}
	}

	return false, nil
}

// Evaluate given user's role by the periodic donation records
func (gs *GormStorage) IsTrailblazer(email string) (bool, error) {
	var result struct {
		Sum int
	}

	query := gs.db.Model(&models.User{}).
		Joins("JOIN periodic_donations p ON users.id = p.user_id").
		Where("p.last_success_at >= DATE_SUB(NOW(), INTERVAL 2 MONTH)").
		Where("users.email = ?", email).
		Group("users.id").
		Select("SUM(p.amount) as sum").
		Scan(&result)

	if query.Error != nil {
		return false, errors.Wrap(query.Error, "failed to check IsTrailblazer")
	}

	log.WithFields(log.Fields{
		"user email": email,
		"sum":        result.Sum,
	}).Info("Checking IsTrailblazer")

	return result.Sum >= 500, nil
}
