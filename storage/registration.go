package storage

import (
	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/models"
)

// RegistrationStorage this is an interface defines methods for users and bookmarks tables
type RegistrationStorage interface {
	// create
	CreateRegistration(models.RegistrationJSON) (models.Registration, error)

	// read
	GetRegistration(string, string) (models.Registration, error)
	GetRegistrationsByService(string, int, int, string, int) ([]models.Registration, error)
	GetRegistrationsAmountByService(string, int) (uint, error)

	// update
	UpdateRegistration(models.RegistrationJSON) (models.Registration, error)

	// delete
	DeleteRegistration(string, string) error
}

// NewGormRegistrationStorage this initializes the user storage
func NewGormRegistrationStorage(db *gorm.DB) RegistrationStorage {
	return &gormRegistrationStorage{db}
}

// gormRegistrationStorage this implements UserStorage interface
type gormRegistrationStorage struct {
	db *gorm.DB
}

func (g *gormRegistrationStorage) GetRegistration(email, service string) (models.Registration, error) {
	var nlr models.Registration
	err := g.db.First(&nlr, "email = ? AND service = ?", email, service).Error
	return nlr, err
}

func (g *gormRegistrationStorage) GetRegistrationsByService(service string, offset, limit int, orderBy string, activeCode int) ([]models.Registration, error) {
	var nlrs []models.Registration

	where := getActiveWhereCondition(activeCode)
	where.Service = service
	err := g.db.Where(&where).Offset(offset).Limit(limit).Order(orderBy).Find(&nlrs).Error
	return nlrs, err
}

func (g *gormRegistrationStorage) GetRegistrationsAmountByService(service string, activeCode int) (uint, error) {
	var count uint

	where := getActiveWhereCondition(activeCode)
	where.Service = service
	err := g.db.Table(constants.RegistrationTable).Where(&where).Count(&count).Error
	return count, err
}

// CreateRegistration this func will create a registration
func (g *gormRegistrationStorage) CreateRegistration(json models.RegistrationJSON) (models.Registration, error) {
	var err error

	registration := models.Registration{
		Email:         json.Email,
		Service:       json.Service,
		Active:        false,
		ActivateToken: json.ActivateToken,
	}

	err = g.db.Create(&registration).Error

	return registration, err
}

// UpdateRegistration this func will update the record in the stroage
func (g *gormRegistrationStorage) UpdateRegistration(json models.RegistrationJSON) (models.Registration, error) {
	var registration models.Registration
	var err error

	err = g.db.First(&registration, "email = ? AND service = ?", json.Email, json.Service).Error

	if err != nil {
		return registration, err
	}

	registration.Email = json.Email
	registration.Service = json.Service
	registration.Active = json.Active
	registration.ActivateToken = json.ActivateToken

	err = g.db.Save(&registration).Error
	return registration, err
}

// DeleteRegistration this func will delete the record in the stroage
func (g *gormRegistrationStorage) DeleteRegistration(email, service string) error {
	err := g.db.Where("email = ? AND service = ?", email, service).Delete(&models.Registration{}).Error
	return err
}

// getActiveWhereCondition recieves 0, 1 or 2.
// 0 means active=false,
// 1 means active=true,
// 2 means active=false || active=true
func getActiveWhereCondition(activeCode int) models.Registration {
	var where models.Registration
	if activeCode == 2 {
		where = models.Registration{}
	} else if activeCode == 1 {
		where = models.Registration{
			Active: true,
		}
	} else {
		where = models.Registration{
			Active: false,
		}
	}

	return where
}
