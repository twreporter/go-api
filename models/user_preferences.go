package models

import (
	"time"
)

type UserPreference struct {
	Topics   []string `json:"topics"`
	Maillist []string `json:"maillist"`
}

type UsersMailGroup struct {
	UserID      int
	MailgroupID string
	CreatedAt   time.Time
}
