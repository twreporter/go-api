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

var InterestIDs = []string{
	"b583571955", // 報導者精選
	"49de5e7889", // 採訪幕後故事
	"092053aaa1", // 報導者營運手記
}
