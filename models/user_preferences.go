package models

import (
	"time"
)

type UserPreference struct {
	ReadPreference []string `json:"read_preference"`
	Maillist       []string `json:"maillist"`
}

type UsersMailGroup struct {
	UserID      int
	MailgroupID string
	CreatedAt   time.Time
}
