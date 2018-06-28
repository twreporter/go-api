package storage

import (
	"fmt"
	"twreporter.org/go-api/models"
)

// CreateAWebPushSubscription - create a record in the persistent database,
// return error if fails.
func (g *GormStorage) CreateAWebPushSubscription(wpSub models.WebPushSubscription) error {
	err := g.db.Create(&wpSub).Error
	if err != nil {
		return g.NewStorageError(err, "GormStorage.CreateAWebPushSubscription", "storage.subscription.error_to_create_a_subscription")
	}

	return nil
}

// GetAWebPushSubscription - read a record from persistent database according to its crc32(endpoint) and endpoint value
func (g *GormStorage) GetAWebPushSubscription(crc32Endpoint uint32, endpoint string) (models.WebPushSubscription, error) {
	var wpSub models.WebPushSubscription
	var err error

	if err = g.db.Find(&wpSub, "crc32_endpoint = ? AND endpoint = ?", crc32Endpoint, endpoint).Error; err != nil {
		return wpSub, g.NewStorageError(err, "GormStorage.GetAWebPushSubscription", fmt.Sprintf("Getting subscription by endpoint %v and crc32_endpoint %v occurs error", endpoint, crc32Endpoint))
	}

	return wpSub, nil
}
