package storage

import (
	"twreporter.org/go-api/constants"
	"twreporter.org/go-api/models"
)

// GetRegistration ...
func (g *GormMembershipStorage) GetRegistration(email, service string) (models.Registration, error) {
	// var svc models.Service
	var reg models.Registration

	svc, err := g.GetServiceByName(service)
	if err != nil {
		return models.Registration{}, err
	}

	err = g.db.Preload("Service", "id = ?", svc.ID).Where("email = ? AND service_id = ?", email, svc.ID).Find(&reg).Error

	return reg, err
}

// GetRegistrationsByService ...
func (g *GormMembershipStorage) GetRegistrationsByService(service string, offset, limit int, orderBy string, activeCode int) ([]models.Registration, error) {
	var regs []models.Registration

	svc, err := g.GetServiceByName(service)
	if err != nil {
		return []models.Registration{}, err
	}

	where := getActiveWhereCondition(activeCode)

	err = g.db.Preload("Service", "id = ?", svc.ID).Where("service_id = ?", svc.ID).Where(where).Offset(offset).Limit(limit).Order(orderBy).Find(&regs).Error
	return regs, err
}

// GetRegistrationsAmountByService ...
func (g *GormMembershipStorage) GetRegistrationsAmountByService(service string, activeCode int) (uint, error) {
	var count uint

	svc, err := g.GetServiceByName(service)
	if err != nil {
		return 0, err
	}

	where := getActiveWhereCondition(activeCode)

	err = g.db.Table(constants.RegistrationTable).Where("service_id = ?", svc.ID).Where(where).Count(&count).Error
	return count, err
}

// CreateRegistration this func will create a registration
func (g *GormMembershipStorage) CreateRegistration(json models.RegistrationJSON) (models.Registration, error) {

	svc, err := g.GetServiceByName(json.Service)
	if err != nil {
		return models.Registration{}, err
	}

	reg := models.Registration{
		ServiceID:     svc.ID,
		Service:       svc,
		Email:         json.Email,
		Active:        false,
		ActivateToken: json.ActivateToken,
	}

	err = g.db.Create(&reg).Error

	return reg, err
}

// UpdateRegistration this func will update the record in the stroage
func (g *GormMembershipStorage) UpdateRegistration(json models.RegistrationJSON) (models.Registration, error) {
	var reg models.Registration

	svc, err := g.GetServiceByName(json.Service)
	if err != nil {
		return models.Registration{}, err
	}

	err = g.db.Where("service_id = ? AND email = ?", svc.ID, json.Email).Find(&reg).Error

	reg.Email = json.Email
	reg.Active = json.Active
	reg.ActivateToken = json.ActivateToken

	err = g.db.Save(&reg).Error
	return reg, err
}

// DeleteRegistration this func will delete the record in the stroage
func (g *GormMembershipStorage) DeleteRegistration(email, service string) error {
	var svc models.Service

	svc, err := g.GetServiceByName(service)
	if err != nil {
		return err
	}

	err = g.db.Where("email = ? AND service_id = ?", email, svc.ID).Delete(&models.Registration{}).Error
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
