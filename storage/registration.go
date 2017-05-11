package storage

import (
	"net/http"

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
	// var svc models.Service
	var reg models.Registration

	// SELECT * FROM `registrations`  WHERE `registrations`.deleted_at IS NULL AND ((email = ${email}))
	// SELECT * FROM `services`  WHERE `services`.deleted_at IS NULL AND (name = ${service})
	err := g.db.Where("email = ?", email).Preload("Service", "name = ?", service).Find(&reg).Error

	return reg, err
}

func (g *gormRegistrationStorage) GetRegistrationsByService(service string, offset, limit int, orderBy string, activeCode int) ([]models.Registration, error) {
	var regs []models.Registration

	where := getActiveWhereCondition(activeCode)

	// SELECT * FROM `registrations`  WHERE `registrations`.deleted_at IS NULL ORDER BY ${orderBy} LIMIT ${limit} OFFSET ${offset}
	// SELECT * FROM `services`  WHERE `services`.deleted_at IS NULL AND (name = ${service})
	err := g.db.Preload("Service", "name = ?", service).Where(where).Offset(offset).Limit(limit).Order(orderBy).Find(&regs).Error
	return regs, err
}

func (g *gormRegistrationStorage) GetRegistrationsAmountByService(service string, activeCode int) (uint, error) {
	var count uint

	where := getActiveWhereCondition(activeCode)

	// SELECT count(*) FROM `registrations`  WHERE (`active` = ${activeCode})
	// SELECT * FROM `services`  WHERE `services`.deleted_at IS NULL AND (name = ${service})
	err := g.db.Table(constants.RegistrationTable).Preload("Service", "name = ?", service).Where(where).Count(&count).Error
	return count, err
}

// CreateRegistration this func will create a registration
func (g *gormRegistrationStorage) CreateRegistration(json models.RegistrationJSON) (models.Registration, error) {
	var err error
	var svc models.Service

	err = g.db.First(&svc, "name = ?", json.Service).Error

	if err != nil {
		return models.Registration{}, models.NewAppError("CreateRegistration",
			"models.registration.create_registration.error_to_get_service", err.Error(), http.StatusInternalServerError)
	}

	reg := models.Registration{
		Service:       svc,
		Email:         json.Email,
		Active:        false,
		ActivateToken: json.ActivateToken,
	}

	err = g.db.Create(&reg).Error

	return reg, err
}

// UpdateRegistration this func will update the record in the stroage
func (g *gormRegistrationStorage) UpdateRegistration(json models.RegistrationJSON) (models.Registration, error) {
	var reg models.Registration

	// SELECT * FROM `registrations`  WHERE `registrations`.deleted_at IS NULL AND ((email = ${email}))
	// SELECT * FROM `services`  WHERE `services`.deleted_at IS NULL AND (name = ${service})
	err := g.db.Where("email = ?", json.Email).Preload("Service", "name = ?", json.Service).Find(&reg).Error

	reg.Email = json.Email
	reg.Active = json.Active
	reg.ActivateToken = json.ActivateToken

	err = g.db.Save(&reg).Error
	return reg, err
}

// DeleteRegistration this func will delete the record in the stroage
func (g *gormRegistrationStorage) DeleteRegistration(email, service string) error {
	var svc models.Service

	g.db.Find(&svc, "name = ?", service)

	err := g.db.Where("email = ? AND service_id = ?", email, svc.ID).Delete(&models.Registration{}).Error
	return err
}

// getActiveWhereCondition recieves 0, 1 or 2.
// 0 means active=false,
// 1 means active=true,
// 2 means active=false || active=true
func getActiveWhereCondition(activeCode int) string {
	var where string
	if activeCode == 2 {
		where = ""
	} else if activeCode == 1 {
		where = "active = 1"
	} else {
		where = "active = 0"
	}

	return where
}
