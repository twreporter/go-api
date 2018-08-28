package utils

import (
	"twreporter.org/go-api/configs/constants"
	"twreporter.org/go-api/models"

	log "github.com/Sirupsen/logrus"
)

// ToNullString invalidates a sql.NullString if empty, validates if not empty
func ToNullString(s string) models.NullString {
	return models.NewNullString(s)
}

// GetNullString returns a invalid NullString
func GetNullString() models.NullString {
	return models.NewNullString("")
}

// GetGender format the gender string
func GetGender(s string) models.NullString {
	var ngender models.NullString
	switch s {
	case "":
		ngender = GetNullString()
	case "male":
		ngender = ToNullString(constants.GenderMale)
	case "female":
		ngender = ToNullString(constants.GenderFemale)
	default:
		// Other gender
		ngender = ToNullString(constants.GenderOthers)
	}
	return ngender
}

// Check - use to fix GoMetaLinter warning of error not check
func Check(f func() error) {
	if err := f(); err != nil {
		log.Error("Received error:", err.Error())
	}
}
