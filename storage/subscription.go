package storage

import (
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

// GetAWebPushSubscriptionByHashEndpoint - read a record from persistent database according to its unique hash endpoint value
func (g *GormStorage) GetAWebPushSubscriptionByHashEndpoint(hashEndpoint string) (models.WebPushSubscription, error) {
	var wpSub models.WebPushSubscription
	var err error

	if err = g.db.Find(&wpSub, "hash_endpoint = ?", hashEndpoint).Error; err != nil {
		return wpSub, g.NewStorageError(err, "GormStorage.GetAWebPushSubscriptionByHashEndpoint", "storage.subscription.error_to_get_a_web_push_subscription")
	}

	return wpSub, nil
}
