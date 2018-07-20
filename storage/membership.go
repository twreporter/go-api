package storage

import (
	"database/sql"

	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/models"
	//log "github.com/Sirupsen/logrus"
)

// MembershipStorage defines the methods we need to implement,
// in order to fulfill the functionalities a membership system needs.
// Such as, let user signup, login w/o oauth, CRUD bookmarks, CRUD registrations.
type MembershipStorage interface {
	/** Close DB Connection **/
	Close() error

	/** User methods **/
	GetUserByID(string) (models.User, error)
	GetUserByEmail(string) (models.User, error)
	GetOAuthData(sql.NullString, string) (models.OAuthAccount, error)
	GetUserDataByOAuth(models.OAuthAccount) (models.User, error)
	GetReporterAccountData(string) (models.ReporterAccount, error)
	GetUserDataByReporterAccount(models.ReporterAccount) (models.User, error)
	InsertOAuthAccount(models.OAuthAccount) error
	InsertReporterAccount(models.ReporterAccount) error
	InsertUserByOAuth(models.OAuthAccount) models.User
	InsertUserByReporterAccount(models.ReporterAccount) (models.User, error)
	UpdateOAuthData(models.OAuthAccount) (models.OAuthAccount, error)
	UpdateReporterAccount(models.ReporterAccount) error

	/** Bookmark methods **/
	GetABookmarkBySlug(string) (models.Bookmark, error)
	GetABookmarkByID(string) (models.Bookmark, error)
	GetABookmarkOfAUser(string, string, string) (models.Bookmark, error)
	GetBookmarksOfAUser(string, int, int) ([]models.Bookmark, int, error)
	CreateABookmarkOfAUser(string, models.Bookmark) (models.Bookmark, error)
	DeleteABookmarkOfAUser(string, string) error

	/** Web Push Subscription methods **/
	CreateAWebPushSubscription(models.WebPushSubscription) error
	GetAWebPushSubscription(uint32, string) (models.WebPushSubscription, error)

	/** Service methods **/
	GetService(string) (models.Service, error)
	CreateService(models.ServiceJSON) (models.Service, error)
	UpdateService(string, models.ServiceJSON) (models.Service, error)
	DeleteService(string) error

	/** Registration methods **/
	GetRegistration(string, string) (models.Registration, error)
	GetRegistrationsByService(string, int, int, string, int) ([]models.Registration, error)
	GetRegistrationsAmountByService(string, int) (uint, error)
	CreateRegistration(string, models.RegistrationJSON) (models.Registration, error)
	UpdateRegistration(string, models.RegistrationJSON) (models.Registration, error)
	DeleteRegistration(string, string) error

	/** Donation methods **/
	CreateAPayByPrimeDonation(models.PayByPrimeDonation) error
	CreateAPeriodDonation(models.PeriodicDonation) error
	CreateAPayByCardTokenDonation(models.PayByCardTokenDonation) error
	CreateAPayByOtherMethodDonation(models.PayByOtherMethodDonation) error
	GetDonationsByPayMethods([]string, uint, uint) (models.DonationRecord, error)
}

// NewGormStorage initializes the storage connected to MySQL database by gorm library
func NewGormStorage(db *gorm.DB) *GormStorage {
	return &GormStorage{db}
}

// GormStorage implements MembershipStorage interface
type GormStorage struct {
	db *gorm.DB
}

// Close quits the DB connection gracefully
func (gs *GormStorage) Close() error {
	err := gs.db.Close()
	if err != nil {
		return err
	}
	return nil
}
