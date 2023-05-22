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

var InterestIDs = map[string]string{
	"featured":            "b583571955", // 報導者精選
	"behind_the_scenes":   "49de5e7889", // 採訪幕後故事
	"operational_journal": "092053aaa1", // 報導者營運手記
}
