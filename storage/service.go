package storage

import (
	"github.com/pkg/errors"

	"github.com/twreporter/go-api/models"
)

// GetService ...
func (g *GormStorage) GetService(name string) (models.Service, error) {
	var s models.Service

	err := g.db.First(&s, "name = ?", name).Error
	if err != nil {
		return s, errors.Wrap(err, "storage.service.get_svc")
	}

	return s, err
}

// CreateService this func will create a service record
func (g *GormStorage) CreateService(json models.ServiceJSON) (models.Service, error) {
	service := models.Service{
		Name: json.Name,
	}

	err := g.db.Create(&service).Error
	if err != nil {
		return service, errors.Wrap(err, "storage.service.create_svc")
	}

	return service, err
}

// UpdateService this func will update the record in the stroage
func (g *GormStorage) UpdateService(name string, json models.ServiceJSON) (models.Service, error) {
	var s models.Service

	err := g.db.Where("name = ?", name).FirstOrCreate(&s).Error
	if err != nil {
		return s, errors.Wrap(err, "storage.service.update_svc")
	}

	s.Name = json.Name

	err = g.db.Save(&s).Error
	if err != nil {
		return s, errors.Wrap(err, "storage.service.update_svc")
	}

	return s, err
}

// DeleteService this func will delete the record in the stroage
func (g *GormStorage) DeleteService(name string) error {
	if g.db.Where("name = ?", name).Delete(&models.Service{}).RowsAffected == 0 {
		return errors.Wrap(ErrRecordNotFound, "storage.service.delete_svc")
	}
	return nil
}
