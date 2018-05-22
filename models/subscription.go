package models

import (
	"time"
)

// TODO add foreign key to bind web push subscription with user later
// WebPushSubscription - a data model which is used by storage to communicate with persistent database
type WebPushSubscription struct {
	ID             uint       `json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
	Endpoint       string     `json:"endpoint"`
	HashEndpoint   string     `gorm:"unique;" json:"hash_endpoint"`
	Keys           string     `json:"-"`
	ExpirationTime *time.Time `json:"expiration_time"`
	UserID         *uint      `json:"user_id"`
}

// SetExpirationTime - set the pointer of expireTime into WebPushSubscription struct.
// The reason why we set the pointer, not the value, into struct
// is because we want to remain the NULL(nil) value if expireTime is not provided.
func (wpSub *WebPushSubscription) SetExpirationTime(expireTime int64) {
	_expireTime := time.Unix(expireTime, 0)
	wpSub.ExpirationTime = &_expireTime
}

// SetUserID - set the pointer of userID into WebPushSubscription struct
func (wpSub *WebPushSubscription) SetUserID(userID uint) {
	wpSub.UserID = &userID
}
