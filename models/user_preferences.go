package models

import (
	"time"
)

type UserPreference struct {
	ReadPreference []string `json:"read_preference"`
	Maillist       []string `json:"maillist"`
}

type UsersMailgroups struct {
	UserID      int
	MailgroupID string
	CreatedAt   time.Time
}
