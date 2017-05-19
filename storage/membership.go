package storage

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/models"
)

// MembershipStorage ...
type MembershipStorage interface {
	/** User methods **/
	GetUserByID(string) (models.User, error)
	GetOAuthData(sql.NullString, string) (models.OAuthAccount, error)
	GetUserDataByOAuth(models.OAuthAccount) (models.User, error)
	GetReporterAccountData(string) (*models.ReporterAccount, error)
	GetUserDataByReporterAccount(*models.ReporterAccount) (*models.User, error)
	InsertUserByOAuth(models.OAuthAccount) models.User
	InsertUserByReporterAccount(models.ReporterAccount) (models.User, error)
	UpdateOAuthData(models.OAuthAccount) (models.OAuthAccount, error)
	UpdateReporterAccountPassword(*models.ReporterAccount, string) (*models.ReporterAccount, error)
	UpdateReporterAccountActive(*models.ReporterAccount, bool) (*models.ReporterAccount, error)

	/** Bookmark methods **/
	GetABookmarkByHref(string) (models.Bookmark, error)
	GetABookmarkByID(string) (models.Bookmark, error)
	GetBookmarksOfAUser(string) ([]models.Bookmark, error)
	CreateABookmarkOfAUser(string, models.Bookmark) error
	DeleteABookmarkOfAUser(string, string) error

	/** Service methods **/
	GetService(string) (models.Service, error)
	CreateService(models.ServiceJSON) (models.Service, error)
	UpdateService(string, models.ServiceJSON) (models.Service, error)
	DeleteService(string) error

	/** Registration methods **/
	GetRegistration(string, string) (models.Registration, error)
	GetRegistrationsByService(string, int, int, string, int) ([]models.Registration, error)
	GetRegistrationsAmountByService(string, int) (uint, error)
	CreateRegistration(models.RegistrationJSON) (models.Registration, error)
	UpdateRegistration(models.RegistrationJSON) (models.Registration, error)
	DeleteRegistration(string, string) error
}

// NewMembershipStorage initializes the storage
func NewMembershipStorage(db *gorm.DB) MembershipStorage {
	return &GormMembershipStorage{db}
}

// GormMembershipStorage implements MembershipStorage interface
type GormMembershipStorage struct {
	db *gorm.DB
}
