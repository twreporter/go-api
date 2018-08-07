package storage

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"twreporter.org/go-api/models"
)

const (
	invalidPeriodicID = 0
)

// CreateAPayByPrimeDonation creates a draft order in database
func (g *GormStorage) CreateAPayByPrimeDonation(m models.PayByPrimeDonation) error {
	errWhere := "GormStorage.CreateAPayByPrimeDonation"
	err := g.db.Create(&m).Error
	if nil != err {
		log.Error(err.Error())
		return g.NewStorageError(err, errWhere, fmt.Sprintf("can not create the record(%#v)", m))
	}
	return nil
}

// UpdateAPayByPrimeDonation updates the draft record with the Tappay response data by order
func (g *GormStorage) UpdateAPayByPrimeDonation(order string, m models.PayByPrimeDonation) error {
	errWhere := "GormStorage.UpdateAPayByPrimeDonation"
	err := g.db.Model(&m).Where("order_number = ?", order).Updates(m).Error
	if nil != err {
		log.Error(err.Error())
		return g.NewStorageError(err, errWhere, fmt.Sprintf("can not update prime donation(order: %s, record: %#v)", order, m))
	}
	return nil
}

// CreateAPeriodicDonation creates the draft record along with the first draft tap pay transaction
func (g *GormStorage) CreateAPeriodicDonation(mpd models.PeriodicDonation, mtd models.PayByCardTokenDonation) (uint, error) {
	errWhere := "GormStorage.CreateAPeriodicDonationWithFirstTransaction"

	tx := g.db.Begin()

	if err := tx.Error; nil != err {
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
		return invalidPeriodicID, g.NewStorageError(err, errWhere, "cannot begin the draft periodic donation creation transaction")
	}

	// Create a draft record for periodic donation
	if err := tx.Create(&mpd).Error; nil != err {
		tx.Rollback()
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
		return invalidPeriodicID, g.NewStorageError(err, errWhere, fmt.Sprintf("cannot create a draft periodic donation(%#v)", mpd))
	}

	// Append the periodic ID for the first record
	mtd.PeriodicID = mpd.ID

	// Create a draft record for the first token transaction
	if err := tx.Create(&mtd).Error; nil != err {
		tx.Rollback()
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
		return invalidPeriodicID, g.NewStorageError(err, errWhere, fmt.Sprintf("cannot create a draft card token donation(%#v)", mtd))
	}

	if err := tx.Commit().Error; nil != err {
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
		return invalidPeriodicID, g.NewStorageError(err, errWhere, "cannot commit the draft periodic donation creation transaction")
	}
	return mpd.ID, nil
}

// DeleteAPeriodicDonation marks the draft tap pay transaction record as 'fail' and then soft delete the periodic donation
func (g *GormStorage) DeleteAPeriodicDonation(periodicID uint, failData models.PayByCardTokenDonation) error {
	errWhere := "GormStorage.DeleteAPeriodicDonation"

	tx := g.db.Begin()

	if err := tx.Error; nil != err {
		return g.NewStorageError(err, errWhere, "cannot begin the draft periodic donation deletion transaction")
	}

	// Updates the draft transaction as fail
	if err := tx.Model(&failData).Where("periodic_id = ?", periodicID).Updates(failData).Error; nil != err {
		tx.Rollback()
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
		return g.NewStorageError(err, errWhere, fmt.Sprintf("cannot update the token record (periodicID: %d, failData: %#v)", periodicID, failData))
	}

	m := models.PeriodicDonation{ID: periodicID}

	// Then, soft-deletes the periodic donation
	if err := tx.Delete(&m).Error; nil != err {
		tx.Rollback()
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
		return g.NewStorageError(err, errWhere, fmt.Sprintf("cannot delete the periodic donation (periodicID: %d)", periodicID))
	}

	if err := tx.Commit().Error; nil != err {
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
		return g.NewStorageError(err, errWhere, "cannot commit the draft periodic donation deletion transaction")
	}
	return nil
}

// UpdateAPeriodicDonation updates the draft tap pay transaction with response and fills up the required information in periodic_donations table
func (g *GormStorage) UpdateAPeriodicDonation(periodicID uint, mpd models.PeriodicDonation, mtd models.PayByCardTokenDonation) error {
	errWhere := "GormStorage.UpdateAPeriodicDonation"
	tx := g.db.Begin()

	if err := tx.Error; nil != err {
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
		return g.NewStorageError(err, errWhere, "cannot begin the draft periodic donation update transaction")
	}
	td := mtd

	// Updates the tap pay success respsonse
	if err := tx.Model(&td).Where("periodic_id = ?", periodicID).Updates(td).Error; nil != err {
		tx.Rollback()
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
		return g.NewStorageError(err, errWhere, fmt.Sprintf("cannot update the token record (data: %#v)", td))
	}

	pd := mpd

	// Then, fills up the required information for later transaction
	if err := tx.Model(&pd).Where("id = ?", periodicID).Updates(pd).Error; nil != err {
		tx.Rollback()
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
		return g.NewStorageError(err, errWhere, fmt.Sprintf("cannot update the periodic donation (data: %#v)", pd))
	}

	if err := tx.Commit().Error; nil != err {
		log.Error(fmt.Sprintf("%s: %s", errWhere, err.Error()))
		return g.NewStorageError(err, errWhere, "cannot commit the draft periodic donation update transaction")
	}
	return nil
}

//TODO
func (g *GormStorage) CreateAPayByOtherMethodDonation(m models.PayByOtherMethodDonation) error {
	return nil
}

//TODO
func (g *GormStorage) GetDonationsByPayMethods(filters []string, offset uint, limit uint) (models.DonationRecord, error) {
	return models.DonationRecord{}, nil
}
