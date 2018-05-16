package models

import (
	"time"
)

// TODO add foreign key to bind webpush subscription with user later
// WebpushSubscription - a data model which is used by storage to communicate with persistent database
type WebpushSubscription struct {
	ID             uint
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
	Endpoint       string
	HashEndpoint   string
	Keys           string
	ExpirationTime *time.Time
	UserID         *uint
}

// SetExpirationTime - set the pointer of expireTime into WebpushSubscription struct.
// The reason why we set the pointer, not the value, into struct
// is because we want to remain the NULL(nil) value if expireTime is not provided.
func (wpSub *WebpushSubscription) SetExpirationTime(expireTime int64) {
	_expireTime := time.Unix(expireTime, 0)
	wpSub.ExpirationTime = &_expireTime
}

// SetUserID - set the pointer of userID into WebpushSubscription struct
func (wpSub *WebpushSubscription) SetUserID(userID uint) {
	wpSub.UserID = &userID
}
