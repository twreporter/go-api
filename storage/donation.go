package storage

import (
	"bytes"

	log "github.com/Sirupsen/logrus"

	"twreporter.org/go-api/models"
)

// CreateAPayByPrimeDonation creates a draft order records in database
func (g *GormStorage) CreateAPayByPrimeDonation(m models.PayByPrimeDonation) error {
	errWhere := "GormStorage.CreateAPayByPrimeDonation"
	err := g.db.Create(&m).Error
	if nil != err {
		log.Error(err.Error())
		return g.NewStorageError(err, errWhere, err.Error())
	}
	return nil
}

// UpdatetTransactionStatus updates the status during the transaction
func (g *GormStorage) UpdateTransactionStatus(order string, status string, table string) error {

	errWhere := "GormStorage.UpdateTransactionStatus"

	var sqlStmt bytes.Buffer
	sqlStmt.WriteString("UPDATE " + table + " ")
	sqlStmt.WriteString("SET status='" + status + "' ")
	sqlStmt.WriteString("WHERE order_number='" + order + "'")

	err := g.db.Exec(sqlStmt.String()).Error
	if nil != err {
		log.Error(err.Error())
		return g.NewStorageError(err, errWhere, err.Error())
	}
	return nil
}

// UpdateAPayByPrimeDonation updates the draft record with the Tappay response data by order
func (g *GormStorage) UpdateAPayByPrimeDonation(order string, m models.PayByPrimeDonation) error {
	errWhere := "GormStorage.UpdateAPayByPrimeDonation"
	err := g.db.Model(&m).Where("order_number = ?", order).Updates(m).Error
	if nil != err {
		log.Error(err.Error())
		return g.NewStorageError(err, errWhere, err.Error())
	}
	return nil
}

//TODO
func (g *GormStorage) CreateAPeriodDonation(m models.PeriodicDonation) error {
	return nil
}

//TODO
func (g *GormStorage) CreateAPayByCardTokenDonation(m models.PayByCardTokenDonation) error {
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
