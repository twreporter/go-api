package storage

import (
	"twreporter.org/go-api/models"
)

// CreateAWebpushSubscription - create a record in the persistent database,
// return error if fails.
func (g *GormStorage) CreateAWebpushSubscription(wpSub models.WebpushSubscription) error {
	err := g.db.Create(&wpSub).Error
	if err != nil {
		return g.NewStorageError(err, "GormStorage.CreateAWebpushSubscription", "")
	}
	return nil
}
