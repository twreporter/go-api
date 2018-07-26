package storage

import (
	"twreporter.org/go-api/models"
)

func (g *GormStorage) CreateAPayByPrimeDonation(m models.PayByPrimeDonation) error {
	errWhere := "GormStorage.CreateAPayByPrimeDonation"
	err := g.db.Create(&m).Error
	if nil != err {
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
