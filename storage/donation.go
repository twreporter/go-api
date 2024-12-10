package storage

import (
	"fmt"
	"time"

	"gopkg.in/guregu/null.v3"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"github.com/twreporter/go-api/models"
)

const (
	YYYYMM      = "200601" //YYYYMM format for time package
	timezoneTPE = "Asia/Taipei"
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

func (g *GormStorage) GetDonationsOfAUser(userID string, limit int, offset int) ([]models.GeneralDonation, int, error) {
	var donations []models.GeneralDonation
	var totalPrime int
	var totalPeriodic int
	var err error

	// build query statement
	defaultColumns := "id, amount, order_number, created_at, send_receipt, status, pay_method, cardholder_first_name, cardholder_last_name, receipt_header, receipt_address_country, receipt_address_state, receipt_address_city, receipt_address_detail, receipt_address_zip_code, card_info_bin_code, card_info_last_four, card_info_type, is_anonymous"
	selectColumnsPrime := fmt.Sprintf("%s, %s", defaultColumns, "'prime' as type")
	queryPrime := g.db.Table("pay_by_prime_donations").Select(selectColumnsPrime).Where("user_id = ?", userID).QueryExpr()
	selectColumnsPeriodic := fmt.Sprintf("%s, %s", defaultColumns, "'periodic' as type")
	queryPeriodic := g.db.Table("periodic_donations").Select(selectColumnsPeriodic).Where("user_id = ?", userID).QueryExpr()
	statement := g.db.Raw("? UNION ?", queryPrime, queryPeriodic).Order("created_at desc").Limit(limit).Offset(offset)

	// build count statement
	primeCountStatement := g.db.Table("pay_by_prime_donations").Where("user_id = ?", userID)
	periodicCountStatement := g.db.Table("periodic_donations").Where("user_id = ?", userID)

	if err = statement.Scan(&donations).Error; err != nil {
		return nil, 0, err
	}
	if err = primeCountStatement.Count(&totalPrime).Error; err != nil {
		return nil, 0, err
	}
	if err = periodicCountStatement.Count(&totalPeriodic).Error; err != nil {
		return nil, 0, err
	}

	return donations, totalPrime + totalPeriodic, nil
}

func (g *GormStorage) GetPaymentsOfAPeriodicDonation(periodicID uint, limit int, offset int) ([]models.Payment, int, error) {
	var payments []models.Payment
	var total int
	var err error

	selectColumns := "created_at, order_number, status, amount"
	statement := g.db.Table("pay_by_card_token_donations").Select(selectColumns).Where("periodic_id = ?", periodicID)
	if err = statement.Order("created_at desc").Limit(limit).Offset(offset).Find(&payments).Error; err != nil {
		return nil, 0, err
	}
	if err = statement.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// TODO
func (g *GormStorage) CreateAPayByOtherMethodDonation(m models.PayByOtherMethodDonation) error {
	return nil
}

func (g *GormStorage) GenerateReceiptSerialNumber(primeID uint, transactionTime null.Time) (string, error) {
	var emptyString string
	needCreateSerial := false
	// get transaction month
	if transactionTime.IsZero() {
		return emptyString, errors.New("transaction time should not be nil")
	}
	tz, err := time.LoadLocation(timezoneTPE)
	if err != nil {
		return emptyString, err
	}
	month := transactionTime.Time.In(tz).Format(YYYYMM)

	// start transaction
	var serialNumber models.ReceiptSerialNumber
	tx := g.db.Begin()

	err = tx.Raw("Select `serial_number` From `receipt_serial_numbers` Where YYYYMM = ? for update", month).Scan(&serialNumber).Error
	if nil != err && !gorm.IsRecordNotFoundError(err) {
		tx.Rollback()
		return emptyString, errors.Wrap(err, fmt.Sprintf("select receipt serial number failed. month: %s", month))
	}
	needCreateSerial = gorm.IsRecordNotFoundError(err)

	// generate receipt number
	receiptNumber := fmt.Sprintf("A%s-%05d", month, serialNumber.SerialNumber+1)

	// update receipt number
	err = tx.Table("pay_by_prime_donations").Where("id = ?", primeID).Update("receipt_number", receiptNumber).Error
	if nil != err {
		tx.Rollback()
		return emptyString, errors.Wrap(err, fmt.Sprintf("update receipt number failed. primeID: %d, receipt number: %s", primeID, receiptNumber))
	}

	// update serial number
	if needCreateSerial {
		serialNumber.YYYYMM = month
		serialNumber.SerialNumber = 1
		if err = tx.Create(&serialNumber).Error; nil != err {
			tx.Rollback()
			return emptyString, errors.Wrap(err, fmt.Sprintf("create receipt serail number failed. YYYYMM: %s", month))
		}
	} else {
		err = tx.Model(&serialNumber).Where("YYYYMM = ?", month).UpdateColumn("serial_number", gorm.Expr("serial_number + ?", 1)).Error
		if nil != err {
			tx.Rollback()
			return emptyString, errors.Wrap(err, fmt.Sprintf("update serial number failed. month: %s, serial number: %d", month, serialNumber.SerialNumber+1))
		}
	}
	tx.Commit()

	return receiptNumber, nil
}
