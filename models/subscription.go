package models

import (
	"time"
)

// TODO add foreign key to bind web push subscription with user later
// WebPushSubscription - a data model which is used by storage to communicate with persistent database
type WebPushSubscription struct {
	ID             uint       `gorm:"primary_key" json:"id"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	Endpoint       string     `gorm:"unique" json:"endpoint"`
	Crc32Endpoint  uint32     `gorm:"index;" json:"crc32_endpoint"`
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

// set WebPushSubscription's table name to be `web_push_subs`
func (WebPushSubscription) TableName() string {
	return "web_push_subs"
}
