package models

import (
	"gopkg.in/guregu/null.v3"
)

type UserPreference struct {
	ReadPreference             []string  `json:"read_preference"`
	ShouldMergeOfflineDonation null.Bool `json:"should_merge_offline_donation_by_identity"`
}
