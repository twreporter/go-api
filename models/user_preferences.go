package models

import (
	"gopkg.in/guregu/null.v3"
)

type UserPreference struct {
	ReadPreference        []string  `json:"read_preference"`
	IsShowOfflineDonation null.Bool `json:"is_showofflinedonation"`
}
