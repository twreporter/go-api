package storage

import (
	"fmt"

	"github.com/pkg/errors"

	"twreporter.org/go-api/models"
)

// CreateAPeriodicDonation creates the draft record along with the first draft tap pay transaction
func (g *GormStorage) CreateAPeriodicDonation(mpd *models.PeriodicDonation, mtd *models.PayByCardTokenDonation) error {
	tx := g.db.Begin()

	if err := tx.Error; nil != err {
		return errors.Wrap(err, "cannot begin the draft periodic donation creation transaction")
	}

	// Create a draft record for periodic donation
	if err := tx.Create(mpd).Error; nil != err {
		tx.Rollback()
		return errors.Wrap(err, fmt.Sprintf("cannot create a draft periodic donation(%#v)", mpd))
	}

	// Append the periodic ID for the first record
	mtd.PeriodicID = mpd.ID

	// Create a draft record for the first token transaction
	if err := tx.Create(mtd).Error; nil != err {
		tx.Rollback()
		return errors.Wrap(err, fmt.Sprintf("cannot create a draft card token donation(%#v)", mtd))
	}

	if err := tx.Commit().Error; nil != err {
		return errors.Wrap(err, "cannot commit the draft periodic donation creation transaction")
	}
	return nil
}

// DeleteAPeriodicDonation marks the draft tap pay transaction record as 'fail' and then soft delete the periodic donation
func (g *GormStorage) DeleteAPeriodicDonation(periodicID uint, failData models.PayByCardTokenDonation) error {

	tx := g.db.Begin()

	if err := tx.Error; nil != err {
		return errors.Wrap(err, "cannot begin the draft periodic donation deletion transaction")
	}

	// Updates the draft transaction as fail
	if err := tx.Model(&failData).Where("periodic_id = ?", periodicID).Updates(failData).Error; nil != err {
		tx.Rollback()
		return errors.Wrap(err, fmt.Sprintf("cannot update the token record (periodicID: %d, failData: %#v)", periodicID, failData))
	}

	m := models.PeriodicDonation{ID: periodicID}

	// Then, soft-deletes the periodic donation
	if err := tx.Delete(&m).Error; nil != err {
		tx.Rollback()
		return errors.Wrap(err, fmt.Sprintf("cannot delete the periodic donation (periodicID: %d)", periodicID))
	}

	if err := tx.Commit().Error; nil != err {
		return errors.Wrap(err, "cannot commit the draft periodic donation deletion transaction")
	}
	return nil
}

// UpdatePeriodicAndCardTokenDonationInTRX updates the draft tap pay transaction with response and fills up the required information in periodic_donations table
func (g *GormStorage) UpdatePeriodicAndCardTokenDonationInTRX(periodicID uint, mpd models.PeriodicDonation, mtd models.PayByCardTokenDonation) error {
	tx := g.db.Begin()

	if err := tx.Error; nil != err {
		return errors.Wrap(err, "cannot begin the draft periodic donation update transaction")
	}
	td := mtd

	// Updates the tap pay success respsonse
	if err := tx.Model(&td).Where("periodic_id = ?", periodicID).Updates(td).Error; nil != err {
		tx.Rollback()
		return errors.Wrap(err, fmt.Sprintf("cannot update the token record (data: %#v)", td))
	}

	pd := mpd

	// Then, fills up the required information for later transaction
	if err := tx.Model(&pd).Where("id = ?", periodicID).Updates(pd).Error; nil != err {
		tx.Rollback()
		return errors.Wrap(err, fmt.Sprintf("cannot update the periodic donation (data: %#v)", pd))
	}

	if err := tx.Commit().Error; nil != err {
		return errors.Wrap(err, "cannot commit the draft periodic donation update transaction")
	}
	return nil
}

//TODO
func (g *GormStorage) CreateAPayByOtherMethodDonation(m models.PayByOtherMethodDonation) error {
	return nil
}
