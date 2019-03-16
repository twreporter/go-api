package storage

import (
	"fmt"

	"twreporter.org/go-api/globals"
	"twreporter.org/go-api/models"
)

// GetRegistration ...
func (g *GormStorage) GetRegistration(email, service string) (models.Registration, error) {
	// var svc models.Service
	var reg models.Registration
	var funcName = "GetRegistration"
	var message = "storage.registration.get_registration"

	svc, err := g.GetService(service)
	if err != nil {
		return reg, err
	}

	err = g.db.Preload("Service", "id = ?", svc.ID).Where("email = ? AND service_id = ?", email, svc.ID).Find(&reg).Error
	if err != nil {
		return reg, g.NewStorageError(err, funcName, fmt.Sprintf("%v.error_to_get_reg: ", message))
	}

	return reg, err
}

// GetRegistrationsByService ...
func (g *GormStorage) GetRegistrationsByService(service string, offset, limit int, orderBy string, activeCode int) ([]models.Registration, error) {
	var regs []models.Registration
	var funcName = "GetRegistrationsByService"
	var message = "storage.registration.get_regs_by_svc"

	svc, err := g.GetService(service)
	if err != nil {
		return regs, err
	}

	where := getActiveWhereCondition(activeCode)

	err = g.db.Preload("Service", "id = ?", svc.ID).Where("service_id = ?", svc.ID).Where(where).Offset(offset).Limit(limit).Order(orderBy).Find(&regs).Error
	if err != nil {
		return regs, g.NewStorageError(err, funcName, fmt.Sprintf("%v.err_to_get_regs: ", message))
	}

	return regs, err
}

// GetRegistrationsAmountByService ...
func (g *GormStorage) GetRegistrationsAmountByService(service string, activeCode int) (uint, error) {
	var count uint
	var funcName = "GetRegistrationsAmountByService"
	var message = "storage.registration.get_regs_amount_by_svc"

	svc, err := g.GetService(service)
	if err != nil {
		return 0, err
	}

	where := getActiveWhereCondition(activeCode)

	err = g.db.Table(globals.RegistrationTable).Where("service_id = ?", svc.ID).Where(where).Count(&count).Error
	if err != nil {
		return 0, g.NewStorageError(err, funcName, fmt.Sprintf("%v.err_to_get_count: ", message))
	}

	return count, err
}

// CreateRegistration this func will create a registration
func (g *GormStorage) CreateRegistration(service string, json models.RegistrationJSON) (models.Registration, error) {
	var user models.User
	var reg models.Registration
	var funcName = "CreateRegistration"
	var message = "storage.registration.create_reg"

	svc, err := g.GetService(service)
	if err != nil {
		return reg, err
	}

	if json.UserID != "" {
		user, err = g.GetUserByID(json.UserID)

		if err != nil {
			return reg, g.NewStorageError(err, funcName, fmt.Sprintf("%v.error_to_get_user", message))
		}
	}

	reg = models.Registration{
		ServiceID:     svc.ID,
		Service:       svc,
		UserID:        user.ID,
		User:          user,
		Email:         json.Email,
		Active:        false,
		ActivateToken: json.ActivateToken,
	}

	err = g.db.Create(&reg).Error
	if err != nil {
		return reg, g.NewStorageError(err, funcName, fmt.Sprintf("%v.error_to_create", message))
	}

	return reg, err
}

// UpdateRegistration this func will update the record in the stroage
func (g *GormStorage) UpdateRegistration(service string, json models.RegistrationJSON) (models.Registration, error) {
	var reg models.Registration
	var funcName = "UpdateRegistration"
	var message = "storage.registration.update_reg"

	svc, err := g.GetService(service)
	if err != nil {
		return reg, err
	}

	err = g.db.Where("service_id = ? AND email = ?", svc.ID, json.Email).FirstOrCreate(&reg).Error
	if err != nil {
		return reg, g.NewStorageError(err, funcName, fmt.Sprintf("%v.error_to_update_reg", message))
	}

	reg.Email = json.Email
	reg.Active = json.Active
	reg.ActivateToken = json.ActivateToken

	err = g.db.Save(&reg).Error
	if err != nil {
		return reg, g.NewStorageError(err, funcName, fmt.Sprintf("%v.error_to_update_reg", message))
	}

	return reg, err
}

// DeleteRegistration this func will delete the record in the stroage
func (g *GormStorage) DeleteRegistration(email, service string) error {
	var svc models.Service
	var funcName = "DeleteRegistration"
	var message = "storage.registration.delete_reg"

	svc, err := g.GetService(service)
	if err != nil {
		return err
	}

	if g.db.Where("email = ? AND service_id = ?", email, svc.ID).Delete(models.Registration{}).RowsAffected == 0 {
		return g.NewStorageError(ErrRecordNotFound, funcName, fmt.Sprintf("%v.error_to_delete", message))
	}

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
