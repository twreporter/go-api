package storage

import (
	"github.com/jinzhu/gorm"
	"twreporter.org/go-api/models"
)

// ServiceStorage this is an interface defines methods for users and bookmarks tables
type ServiceStorage interface {
	// create
	CreateService(models.ServiceJSON) (models.Service, error)

	// read
	GetService(string) (models.Service, error)

	// update
	UpdateService(string, models.ServiceJSON) (models.Service, error)

	// delete
	DeleteService(string) error
}

// NewGormServiceStorage this initializes the user storage
func NewGormServiceStorage(db *gorm.DB) ServiceStorage {
	return &gormServiceStorage{db}
}

// gormServiceStorage this implements UserStorage interface
type gormServiceStorage struct {
	db *gorm.DB
}

func (g *gormServiceStorage) GetService(id string) (models.Service, error) {
	var s models.Service
	err := g.db.First(&s, "id = ?", id).Error
	return s, err
}

func (g *gormServiceStorage) GetServiceByName(name string) (models.Service, error) {
	var s models.Service
	err := g.db.First(&s, "name = ?", name).Error
	return s, err
}

// CreateService this func will create a service record
func (g *gormServiceStorage) CreateService(json models.ServiceJSON) (models.Service, error) {
	var err error

	service := models.Service{
		Name: json.Name,
	}

	err = g.db.Create(&service).Error

	return service, err
}

// UpdateService this func will update the record in the stroage
func (g *gormServiceStorage) UpdateService(id string, json models.ServiceJSON) (models.Service, error) {
	var s models.Service
	var err error

	err = g.db.First(&s, "id = ?", id).Error

	if err != nil {
		return s, err
	}

	s.Name = json.Name
	err = g.db.Save(&s).Error
	return s, err
}

// DeleteService this func will delete the record in the stroage
func (g *gormServiceStorage) DeleteService(id string) error {
	err := g.db.Where("id = ?", id).Delete(&models.Service{}).Error
	return err
}
