package storage

import (
	"twreporter.org/go-api/models"
)

//TODO
func (g *GormStorage) CreateAPayByPrimeDonation(m models.PayByPrimeDonation) error {
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
