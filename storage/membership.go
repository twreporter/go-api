package storage

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/jinzhu/gorm"
	"gopkg.in/guregu/null.v3"

	"twreporter.org/go-api/models"
)

// MembershipStorage defines the methods we need to implement,
// in order to fulfill the functionalities a membership system needs.
// Such as, let user signup, login w/o oauth, CRUD bookmarks, CRUD registrations.
type MembershipStorage interface {
	/** Close DB Connection **/
	Close() error

	/** Default CRUD **/
	Create(interface{}) error
	Get(uint, interface{}) error
	GetByConditions(map[string]interface{}, interface{}) error
	Update(interface{}) error
	UpdateByConditions(map[string]interface{}, interface{}) error
	Delete(uint, interface{}) error

	/** User methods **/
	GetUserByID(string) (models.User, error)
	GetUserByEmail(string) (models.User, error)
	GetOAuthData(null.String, string) (models.OAuthAccount, error)
	GetUserDataByOAuth(models.OAuthAccount) (models.User, error)
	GetReporterAccountData(string) (models.ReporterAccount, error)
	GetUserDataByReporterAccount(models.ReporterAccount) (models.User, error)
	InsertOAuthAccount(models.OAuthAccount) error
	InsertReporterAccount(models.ReporterAccount) error
	InsertUserByOAuth(models.OAuthAccount) (models.User, error)
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

	/** Donation methods **/
	CreateAPeriodicDonation(models.PeriodicDonation, models.PayByCardTokenDonation) (uint, error)
	DeleteAPeriodicDonation(uint, models.PayByCardTokenDonation) error
	UpdatePeriodicAndCardTokenDonationInTRX(uint, models.PeriodicDonation, models.PayByCardTokenDonation) error
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

// Get method of MembershipStorage interface
func (gs *GormStorage) Get(id uint, m interface{}) error {
	var err error
	var errWhere string = "GormStorage.Get"

	err = gs.db.Where("id = ?", id).Find(m).Error

	if err != nil {
		log.Error(err.Error())
		return gs.NewStorageError(err, errWhere, fmt.Sprintf("can not get the record(id: %d)", id))
	}

	return nil
}

// GetByConditions method of MembershipStorage interface
func (gs *GormStorage) GetByConditions(cond map[string]interface{}, m interface{}) error {
	var err error
	var errWhere string = "GormStorage.GetByConditions"

	err = gs.db.Where(cond).Find(m).Error

	if err != nil {
		log.Error(err.Error())
		return gs.NewStorageError(err, errWhere, fmt.Sprintf("can not get the record(where: %v)", cond))
	}

	return nil
}

// Update method of MembershipStorage interface
func (gs *GormStorage) Update(m interface{}) error {
	var err error
	var errWhere string = "GormStorage.Update"

	// caution:
	// batch updates if primary key of m is zero value
	err = gs.db.Model(m).Updates(m).Error

	if err != nil {
		log.Error(err.Error())
		return gs.NewStorageError(err, errWhere, fmt.Sprintf("can not update the record(%#v)", m))
	}

	return nil
}

// UpdateByConditions method of MembershipStorage interface
func (gs *GormStorage) UpdateByConditions(cond map[string]interface{}, m interface{}) error {
	var err error
	var errWhere string = "GormStorage.UpdateByConditions"

	// caution:
	// batch updates if cond is zero value and primary key of m is zero value
	err = gs.db.Model(m).Where(cond).Updates(m).Error

	if err != nil {
		log.Error(err.Error())
		return gs.NewStorageError(err, errWhere, fmt.Sprintf("can not update the record(where: %v)", cond))
	}

	return nil
}

// Delete method of MembershipStorage interface
func (gs *GormStorage) Delete(id uint, m interface{}) error {
	return nil
}

// Create method of MembershipStorage interface
func (gs *GormStorage) Create(m interface{}) error {
	var err error
	var errWhere string = "GormStorage.Create"

	err = gs.db.Create(m).Error

	if nil != err {
		log.Error(err.Error())
		return gs.NewStorageError(err, errWhere, fmt.Sprintf("can not create the record(%#v)", m))
	}

	return nil
}
