package models

import (
	"time"

	"gopkg.in/guregu/null.v3"
)

type UserPreference struct {
	ReadPreference        []string  `json:"read_preference"`
	Maillist              []string  `json:"maillist"`
	IsShowOfflineDonation null.Bool `json:"is_showofflinedonation"`
}

type UsersMailgroups struct {
	UserID      int
	MailgroupID string
	CreatedAt   time.Time
}
