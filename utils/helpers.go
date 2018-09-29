package utils

import (
	"gopkg.in/guregu/null.v3"
	"twreporter.org/go-api/configs/constants"

	log "github.com/Sirupsen/logrus"
)

// GetGender format the gender string
func GetGender(s string) null.String {
	var gender string
	switch s {
	case "":
		gender = s
	case "male":
		gender = constants.GenderMale
	case "female":
		gender = constants.GenderFemale
	default:
		// Other gender
		gender = constants.GenderOthers
	}
	return null.StringFrom(gender)
}

// Check - use to fix GoMetaLinter warning of error not check
func Check(f func() error) {
	if err := f(); err != nil {
		log.Error("Received error:", err.Error())
	}
}
